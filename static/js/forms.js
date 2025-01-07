// Remove a parent block
function removeParent(button) {
    button.parentElement.remove();
}

function removeTarget(button) {
    $(button.getAttribute("data-target")).remove();
    // button.parentElement.remove();
}

function removeLastOf(button) {
    $(`${button.getAttribute("data-target")}:last-of-type`).remove();
    // button.parentElement.remove();
}
