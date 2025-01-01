const displayFormGroup = (serviceCount, id, labelText, type, placeholder) => {
    const newNode = $("#shadow-form .form-group").cloneNode(true);
    const label = newNode.querySelector("label");
    label.setAttribute("for", `services[${serviceCount}][${id}]`);
    label.innerHTML = labelText;
    const input = newNode.querySelector("input");
    input.setAttribute("id", `services[${serviceCount}][${id}]`);
    input.setAttribute("type", type);
    input.setAttribute("name", `services[${serviceCount}][${id}]`);
    input.setAttribute("placeholder", placeholder);
    return newNode;
}
let serviceCount = document.querySelectorAll("#servicesContainer .service-block").length; // Track the number of service blocks

// Add a new service block
try {
    document.getElementById('addServiceButton').addEventListener('click', function () {
        const container = $('#servicesContainer');
        const newServiceBlock = document.createElement('div');
        newServiceBlock.className = 'service-block bg-slate-800 border-gray-900 border-1 rounded p-4 mb-2.5';
        try {
            newServiceBlock.appendChild(displayFormGroup(serviceCount, "id", "Name", "text", "Service name"));
            newServiceBlock.appendChild(displayFormGroup(serviceCount, "method", "Method", "text", "string, regex etc..."));
            newServiceBlock.appendChild(displayFormGroup(serviceCount, "pattern", "Pattern", "text", "/api/vault_hunter"));
            newServiceBlock.appendChild(displayFormGroup(serviceCount, "port", "Port", "number", "Port of the service"));
            newServiceBlock.appendChild($("#shadow-form .remove-button").cloneNode(true));
        } catch (e) {
            ToastError("ERROR: Could not display a new form group")
            return;
        }
        container.appendChild(newServiceBlock);
        serviceCount++;
    });
} catch (e) {
    console.log(e)
} 

// Remove a parent block
function removeParent(button) {
    button.parentElement.remove();
}

function makeServiceLabel(it, name) {
    return `services[${it}][${name}]`;
}