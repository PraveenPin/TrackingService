cd ~/Downloads/go/

scp -r -i TrackingService/ec2keypair.pem TrackingService.zip ec2-user@ec2-18-191-253-123.us-east-2.compute.amazonaws.com:/home/ec2-user/

ssh -i TrackingService/ec2keypair.pem ec2-user@ec2-18-191-253-123.us-east-2.compute.amazonaws.com

cd TrackingService
go run main.go

#install gcloud
curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-429.0.0-linux-x86_64.tar.gz

tar -xf google-cloud-cli-429.0.0-linux-x86_64.tar.gz

./google-cloud-sdk/install.sh

./google-cloud-sdk/bin/gcloud init

gcloud auth application-default login