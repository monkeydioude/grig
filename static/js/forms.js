function updateFormGroup(node, serviceCount, id) {
    node.querySelector("label").setAttribute("for", `services[${serviceCount}][${id}]`);
    const input = node.querySelector("input");
    input.setAttribute("id", `services[${serviceCount}][${id}]`);
    input.setAttribute("name", `services[${serviceCount}][${id}]`);
}

// addFormGroupListener listens to the click event on a chosen node,
// spawns inside a new block of service blocks
// childsList represents the list of child input, in order,
// that should be modified.
function addFormGroupListener(counter, listenOn, container, modelNode, childsList) {
    try {
        listenOn.addEventListener('click', function () {
            try {
                const newBlock = modelNode.querySelector(".service-block").cloneNode(true);
                for (const it in childsList) {
                    updateFormGroup(newBlock.childNodes[it], counter, childsList[it]);
                }
                newBlock.appendChild(modelNode.querySelector(".remove-button").cloneNode(true));
                container.appendChild(newBlock);
                counter++;
            } catch (e) {
                throw e;
            }
        });
    } catch (e) {
        console.log(e)
        ToastError("ERROR: Could not display a new form group")
    }
}

// Remove a parent block
function removeParent(button) {
    button.parentElement.remove();
}