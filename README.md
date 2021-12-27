# twitch-bot

Build
```
go install github.com/google/ko@latest
```

```
gcloud auth configure-docker -q
IMAGE="$(KO_DOCKER_REPO=gcr.io/ahmet-personal-api/twitch-bot ko publish .)"
```

First-time set up

```
TOKEN=xxx
```

```
gcloud compute instances create-with-container -q --project ahmet-personal-api \
  --machine-type e2-micro --zone us-west1-b \
  --container-image="$IMAGE" \
  --container-restart-policy=always \
  --container-env=TWITCH_USER=Ahmet_Alp \
  --container-env=TWITCH_TOKEN="$TOKEN" \
  free-vm
```

Update:

```
gcloud compute instances update-container -q --project ahmet-personal-api \
  --zone us-west1-b \
  --container-image="$IMAGE" \
  free-vm
```

Logs:

```
gcloud compute ssh -q --project ahmet-personal-api --zone us-west1-b free-vm
docker logs -f [...]
```
