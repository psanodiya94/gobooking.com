
function Prompt() {
    let toast = function(options) {
        const {
            title = "",
            icon = "success",
            position = "top-end",
        } = options;

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });

        Toast.fire({ });
    }

    let success = function(options) {
        const {
            title = "",
            text = "",
            footer = "",
        } = options;

        Swal.fire({
            icon: "success",
            title: title,
            text: text,
            footer: footer
        });
    }

    let error = function(options) {
        const {
            title = "",
            text = "",
            footer = "",
        } = options;

        Swal.fire({
            icon: "error",
            title: title,
            text: text,
            footer: footer
        });
    }

    async function custom(options) {
        const {
            icon = "",
            title = "",
            text = "",
            showCancelButton = true,
            showConfirmButton = true,
        } = options;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: text,
            focusConfirm: false,
            showCancelButton: showCancelButton,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (options.willOpen !== undefined) {
                    options.willOpen();
                }
            },
            didOpen: () => {
                if (options.didOpen !== undefined) {
                    options.didOpen();
                }
            },
        });

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (options.callback !== undefined) {
                        options.callback(result);
                    }
                } else {
                    options.callback(false);
                }
            } else {
                options.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}