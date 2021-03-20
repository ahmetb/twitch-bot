package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v2"
)

type Msg = twitch.PrivateMessage

const (
	channel = `hasanabi`
	minWait = time.Second*2 + time.Millisecond*500
	maxWait = time.Second * 10
)

var (
	re1 = regexp.MustCompile(`(?i).*t[uü]rk\s+var\s*m[ıi].*`)
	re2 = regexp.MustCompile(`(?i)^(@hasanabi|abi)?\s*t[uü]rk\s*m[uü]s[uü]n.*`)
	re3 = regexp.MustCompile(`(?i).*ahmet_alp.*`)

	minPerUserReactInterval = time.Minute

	mu         sync.RWMutex
	userReacts = make(map[string]time.Time)
)

func main() {
	client := twitch.NewClient(os.Getenv(`TWITCH_USER`), os.Getenv("TWITCH_TOKEN"))
	client.Join(channel)
	client.OnConnect(func() {
		client.Say(channel, "hello everyone!"+mkFlag(5))
		fmt.Println("connected.")
	})
	msgs := make(chan Msg)
	client.OnPrivateMessage(onMsg(msgs))

	ch1 := make(chan Msg)
	ch2 := make(chan Msg)
	ch3 := make(chan Msg)

	go func() {
		for m := range msgs {
			if re1.MatchString(m.Message) {
				ch1 <- m
			} else if re2.MatchString(m.Message) {
				ch2 <- m
			} else if re3.MatchString(m.Message) {
				ch3 <- m
			}
		}
	}()
	go consume(client, delayed(ch1), turkVarMiReact)
	go consume(client, delayed(ch2), abiTurkMusunReact)
	go printMention(ch3)

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func delayed(c chan Msg) chan Msg {
	outCh := make(chan Msg)
	go func() {
		for m := range c {
			go func(v Msg) {
				time.Sleep(
					minWait + time.Duration(rand.Int63n(int64(maxWait-minWait))))
				outCh <- v
			}(m)
		}
	}()
	return outCh
}

func onMsg(c chan Msg) func(Msg) {
	return func(msg Msg) {
		c <- msg
	}
}

func consume(client *twitch.Client, c <-chan Msg, react func() string) {
	for m := range c {
		user := m.User.DisplayName

		fmt.Printf("%s %s",
			color.RedString("@"+user),
			color.HiWhiteString(m.Message),
		)

		mu.RLock()
		v := time.Since(userReacts[user])
		mu.RUnlock()
		if v < minPerUserReactInterval {
			fmt.Println("-->" + color.HiRedString("THROTTLED (%v)", v.Truncate(time.Second)))
			continue
		}

		resp := fmt.Sprintf(`@%s %s`, user, react())
		client.Say(channel, resp)

		mu.Lock()
		userReacts[user] = time.Now()
		mu.Unlock()
		fmt.Println(" --> " + color.HiBlackString(resp))
	}
}

func turkVarMiReact() string {
	v := []string{
		`var olm yeter sormayın artık`,
		`var türk mnkym var`,
		`var olm sorup durmayın` + mkFlag(2),
		`var niye lazim sana türk bulup napıcan ` + mkFlag(2),
		`türk var kardeşim, nargile mi içmek istiyorsun beraber ` + mkFlag(3),
		`türk var kardeşim, madem türk'sün göster ürksün ` + mkFlag(5),
		`var kardeşi̇m otağı nereye kuruyoruz yer göster ` + mkFlag(5),
		mkFlag(5) + ` var kardeşim, VER MEHTERİ CcC ` + mkFlag(5),
		`türk olmaz olur mu kardeşim her yerde türk var!` + mkFlag(1),
		`bütün dünya dolaylı olarak türk olduğu için TÜRK VAR evet`,
		`bütün dünya dolaylı olarak türk` + mkFlag(1) + ` olduğu için TÜRK VAR evet`,
		`türk var hatta GÖKTÜRK bile var lan cCc` + mkFlag(1),
		`"türk varmı diye sorulmadığı gün muhasız medeniyetler seviyesine erişmiş olacağız." --K.Atatürk`,
		`"Beni chatte TÜRK VARMI diye soranlara emanet ediniz." - Kemal Ataturk`,
		`Türk yok kardeşim burda bütün Türkiye'ye nükleer bomba atılmış bi Türk sen kalmışsın.` + mkFlag(2),
		`Türk yok ama GÖK TÜRK var, o olur mu?` + mkFlag(4),
		`Türk yok ama ANADOLU SELÇUKLU var, o kurtarir mi?`,
	}
	return v[rand.Intn(len((v)))]
}

func abiTurkMusunReact() string {
	v := []string{
		mkFlag(2) + `evet türk olm sorup durmayın` + mkFlag(2),
		`evet türk ne yapacan ` + mkFlag(2),
		`evet Hasan türk` + mkFlag(1) + ` ama asıl soru sen Türk müsün? Eğer öyleysen say bütün vergileri.`,
		`millete türkmüsün diye soracağına sen türk müsün asıl`,
		`Kendine türk'üm diyen türk'tür. -K.Atatürk` + mkFlag(4),
		`evet türk ne yapacan, nikah mı kıycan `,
		`evet adam türk ne yapacan, nüfüsüna mı geçiricen?`,
	}
	return v[rand.Intn(len((v)))]
}

func mkFlag(n int) string {
	return strings.Repeat(`🇹🇷`, rand.Intn(n))
}

func printMention(c <-chan Msg) {
	for m := range c {
		fmt.Printf("%s %s\n", color.GreenString("@"+m.User.DisplayName),
			color.YellowString(m.Message))
	}
}
