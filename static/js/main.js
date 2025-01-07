// THE FEELS! Old jQuery like selector shortcuts!! We went full circle
const $ = function(selectorOrNode, thenSelector) {
    if (selectorOrNode instanceof Node) {
        return selectorOrNode.querySelector(thenSelector);
    }
    return document.querySelector(selectorOrNode);
}

const $$ = function(selectorOrNode, thenSelector) {
    if (selectorOrNode instanceof Node) {
        return selectorOrNode.querySelectorAll(thenSelector);
    }
    return document.querySelectorAll(selectorOrNode);
};

const $_ = function(node, thenSelector) {
    const directChildren = Array.from(node.children);
    return directChildren.filter(child => child.classList.contains(thenSelector));
}

if (!Node.prototype.removeClass) {
    Node.prototype.removeClass = function (className) {
        if (this.classList) {
            this.classList.remove(className);
        } else if (this.className) {
            // Fallback for older browsers
            this.className = this.className
                .split(' ')
                .filter(c => c !== className)
                .join(' ');
        }
    };
}

if (!Node.prototype.addClass) {
    Node.prototype.addClass = function (className) {
        if (this.classList) {
            // Use classList if supported
            this.classList.add(className);
        } else if (this.className !== undefined) {
            // Fallback for older browsers
            const classes = this.className.split(' ');
            if (!classes.includes(className)) {
                classes.push(className);
                this.className = classes.join(' ');
            }
        }
    };
}

function Toast(text, style, className) {
    Toastify({
        text: text,
        duration: 3000,
        gravity: "top", // `top` or `bottom`
        position: "center", // `left`, `center` or `right`
        stopOnFocus: true, // Prevents dismissing of toast on hover
        className: `rounded ${className || ""}`,
        style,
        onClick: function () { } // Callback after click
    }).showToast();
}

function ToastSuccess(text) {
    Toast(text, {
        background: "linear-gradient(to right, #00b09b, #96c93d)",
    });
}

function ToastError(text) {
    Toast(text, {
        background: "linear-gradient(to right, #ff6347, #cc0000)",
    });
}