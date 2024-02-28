
const host = "https://pkg.go.dev"
const moduleName = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v2@v2.0.0";

(async () => {
    const path = `${host}/fetch/${moduleName}`
    console.log("path:", path)
    const response = await fetch(`${host}/fetch/${moduleName}`, { method: 'POST' });

    if (response.ok) {
        console.log("success!!!")
        return;
    }

    console.log(await response.text());
})();

