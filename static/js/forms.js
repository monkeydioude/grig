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

// --- Service modals -------------------------------------------------
// Each service row opens a <dialog> holding its fields. The fields stay
// inside the page's form the whole time, so a closed dialog still
// submits and no value has to be copied back and forth.

function openDialog(trigger) {
    const dialog = document.getElementById(trigger.dataset.dialog);
    if (!dialog) {
        return;
    }
    // only one at a time, otherwise dialogs stack on top of each other
    document.querySelectorAll("dialog[open]").forEach(open => {
        if (open !== dialog) {
            open.close();
        }
    });
    dialog.showModal();
}

function closeDialog(element) {
    const dialog = element.closest("dialog");
    if (dialog) {
        dialog.close();
    }
}

// renumberServices rewrites the services[N][...] field names so the
// posted array has no gaps after a row is removed, and no duplicate
// index after one is appended.
function renumberServices(container) {
    Array.from(container.children)
        .filter(child => child.classList.contains("service-block"))
        .forEach((block, index) => {
            block.querySelectorAll("[name]").forEach(field => {
                field.name = field.name.replace(/^services\[\d+\]/, `services[${index}]`);
            });
        });
}

// htmx:afterSwap bubbles, so a swap anywhere inside the services list --
// including the key/value rows appended inside a modal -- reaches the
// list's handler. Only the list itself gaining a row should react.
function onServicesSwap(event, container) {
    if (event.target !== container) {
        return;
    }
    renumberServices(container);
    openLastServiceDialog(container);
}

// A new row is appended empty, so open it right away rather than leaving
// a blank line in the list.
function openLastServiceDialog(container) {
    const blocks = Array.from(container.children).filter(c => c.classList.contains("service-block"));
    const last = blocks[blocks.length - 1];
    const trigger = last && last.querySelector("[data-dialog]");
    if (trigger) {
        openDialog(trigger);
    }
}

let pendingServiceDelete = null;

function askDeleteService(button) {
    pendingServiceDelete = button.closest(".service-block");
    closeDialog(button);
    const name = pendingServiceDelete.querySelector('[name$="[id]"]');
    $("#confirm-delete-name").textContent = (name && name.value) || "this service";
    document.getElementById("confirm-delete-service").showModal();
}

function confirmDeleteService() {
    document.getElementById("confirm-delete-service").close();
    if (!pendingServiceDelete) {
        return;
    }
    const form = pendingServiceDelete.closest("form");
    const container = pendingServiceDelete.parentElement;
    pendingServiceDelete.remove();
    pendingServiceDelete = null;
    renumberServices(container);
    // requestSubmit fires a real submit event, which htmx intercepts
    form.requestSubmit();
}

// Close any open dialog once a save round-trips successfully. Only a
// form submit counts: the buttons that append a row inside a modal are
// htmx requests too, and closing on those would shut the modal the user
// is still working in.
document.addEventListener("htmx:afterRequest", function (event) {
    const source = event.detail.elt;
    const saved = event.detail.xhr && event.detail.xhr.status === 200;
    if (saved && source && source.tagName === "FORM") {
        // hidden extra-field rows have already submitted their empty
        // value to delete the key on disk — remove them so the
        // duplicate-key check in bindExtraKey does not block re-use.
        source.querySelectorAll(".extra-field.hidden").forEach(row => row.remove());
        document.querySelectorAll("dialog[open]").forEach(dialog => dialog.close());
    }
});

// --- Extra keys -----------------------------------------------------
// Keys grig does not model are edited as free key/value pairs. Field
// names are derived from the service the row sits in, so they survive
// the renumbering that a delete triggers.

const RESERVED_SERVICE_KEYS = ["id", "method", "pattern", "host", "port", "protocol"];

function servicePrefix(element) {
    const idField = element.closest(".service-block").querySelector('input[name$="[id]"]');
    return idField ? idField.name.replace(/\[id\]$/, "") : null;
}

// bindExtraKey names the value input after the key the user typed. An
// empty, reserved or already-used key leaves the input unnamed, which
// means it is simply not submitted.
function bindExtraKey(input) {
    const row = input.closest(".extra-field");
    const valueField = row.querySelector(".extra-value");
    const block = input.closest(".service-block");
    valueField.removeAttribute("name");

    // also remove any previous type hint
    const prevHint = row.querySelector(".extra-type-hint");
    if (prevHint) prevHint.removeAttribute("name");

    const key = input.value.trim().replace(/[\[\]]/g, "");
    let problem = "";
    if (key !== "") {
        if (RESERVED_SERVICE_KEYS.includes(key)) {
            problem = `"${key}" is already a field above`;
        } else if (block.querySelector(`[name$="[${key}]"]`)) {
            problem = `"${key}" is already set`;
        }
    }
    row.querySelector(".extra-error").textContent = problem;
    input.classList.toggle("border-red-500", problem !== "");

    const prefix = servicePrefix(input);
    if (key !== "" && problem === "" && prefix) {
        valueField.name = `${prefix}[${key}]`;
        // set the type hint companion field so the backend knows
        // the intended type for brand-new keys
        const typeSelect = row.querySelector(".extra-type-select");
        if (typeSelect) {
            let hint = row.querySelector(".extra-type-hint");
            if (!hint) {
                hint = document.createElement("input");
                hint.type = "hidden";
                hint.classList.add("extra-type-hint");
                row.appendChild(hint);
            }
            hint.name = `${prefix}[__type__${key}]`;
            hint.value = typeSelect.value;
        }
    }
}

// changeExtraType swaps the value widget based on the type dropdown:
//   string → <input type="text">
//   number → <input type="number">
//   bool   → <select> with true/false
function changeExtraType(select) {
    const row = select.closest(".extra-field");
    const container = row.querySelector(".extra-value-container");
    const oldField = container.querySelector(".extra-value");
    const oldName = oldField ? oldField.getAttribute("name") : null;
    const type = select.value;

    let newField;
    if (type === "bool") {
        newField = document.createElement("select");
        newField.innerHTML = '<option value="true">true</option><option value="false">false</option>';
    } else {
        newField = document.createElement("input");
        newField.type = type === "number" ? "number" : "text";
        newField.placeholder = type === "number" ? "0" : "/somewhere";
    }
    // preserve shared classes: form input styling + marker class
    newField.className = oldField.className;
    if (oldName) newField.name = oldName;

    oldField.replaceWith(newField);

    // also update the type hint value if it exists
    const hint = row.querySelector(".extra-type-hint");
    if (hint) hint.value = type;
}

// removeExtraField empties a saved key rather than detaching its row: a
// key that is not submitted at all is restored from the file on save, so
// posting it empty is what actually deletes it. The original value is
// kept in data-original-value so a cancel can undo the removal.
function removeExtraField(button) {
    const row = button.closest(".extra-field");
    const valueField = row.querySelector(".extra-value");
    if (valueField && valueField.hasAttribute("name")) {
        row.dataset.originalValue = valueField.value;
        valueField.value = "";
        row.classList.add("hidden");
    } else {
        row.remove();
    }
}

// When a dialog is closed without saving (Cancel or Escape), restore
// any extra-field rows that were hidden by removeExtraField and discard
// unsaved new rows added via "+ Add key". After a successful save the
// htmx:afterRequest handler already removes the hidden rows, so the
// close listener finds nothing to restore.
document.addEventListener("close", function (event) {
    if (event.target.tagName !== "DIALOG") return;
    const dialog = event.target;

    // restore deleted-but-not-saved existing rows
    dialog.querySelectorAll(".extra-field.hidden").forEach(row => {
        const valueField = row.querySelector(".extra-value");
        if (valueField && row.dataset.originalValue !== undefined) {
            valueField.value = row.dataset.originalValue;
            delete row.dataset.originalValue;
        }
        row.classList.remove("hidden");
    });

    // remove unsaved new rows (they carry the type selector)
    dialog.querySelectorAll(".extra-field .extra-type-select").forEach(sel => {
        sel.closest(".extra-field").remove();
    });
}, true);
