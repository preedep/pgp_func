rm -rf deployment.zip

echo "building linux bin"
GOOS=linux GOARCH=amd64 go build -o handler main.go
if [ "$?" != "0" ]; then
echo "[Error] building!" 1>&2
exit 1
fi
####
echo "ziping deployment package"
zip -r deployment.zip handler host.json HttpTriggerPGP/ EventGridTriggerBlobCreated/ BlobTriggerPGP/
if [ "$?" != "0" ]; then
echo "[Error] ziping!" 1>&2
exit 1
fi
#####
echo "deploying to azure function"
az functionapp deployment source config-zip \
-g mydevresourcesg101 -n pgp-func \
--src ./deployment.zip
if [ "$?" != "0" ]; then
echo "[Error] deploying!" 1>&2
exit 1
fi
#####

echo "complete"

