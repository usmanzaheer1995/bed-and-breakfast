function Prompt() {
    const toast = function (c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        });

        Toast.fire();
    }

    const success = function (c) {
        const {
            msg = '',
            title = '',
            footer = '',
        } = c;

        Swal.fire({
            icon: 'success',
            title,
            text: msg,
            footer,
        })
    }

    const error = function (c) {
        const {
            msg = '',
            title = '',
            footer = '',
        } = c;

        Swal.fire({
            icon: 'error',
            title,
            text: msg,
            footer,
        })
    }

    async function custom(c) {
        const {
            icon = '',
            msg = '',
            title = '',
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon,
            title,
            html: msg,
            backdrop: false,
            allowOutsideClick: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton,
            willOpen() {
                if (c.willOpen) c.willOpen();
            },
            didOpen() {
                if (c.didOpen) c.didOpen();
            },
        });

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result === true) {
                    if (c.callback !== undefined) c.callback(result);
                    return;
                }
                if (!result.includes("")) {
                    if (c.callback !== undefined) c.callback(result);
                }
            } else {
                c.callback(false);
            }
        } else {
            c.callback(false);
        }
    }

    return {
        toast,
        success,
        error,
        custom,
    };
}

function checkRoomAvailability(room_id, csrf_token) {
    const attention = Prompt()
    const html = `
        <form id="check-availability-form" action-="" autocomplete="off" method="post" novalidate class="needs-validation">
          <div class="row">
            <div class="col">
              <div class="row" id="reservation-dates-modal">
                <div class="col">
                  <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                </div>
                <div class="col">
                  <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                </div>
              </div>
            </div>
          </div>

        </form>
  `;
    attention.custom({
        msg: html,
        title: 'Choose your dates',
        willOpen() {
            const elem = document.getElementById("reservation-dates-modal");
            const rp = new DateRangePicker(elem, {
                format: "yyyy-mm-dd",
                showOnFocus: true,
                minDate: new Date(),
            });
        },
        didOpen() {
            document.getElementById('start').removeAttribute('disabled');
            document.getElementById('end').removeAttribute('disabled');
        },
        callback(result) {
            if (result) {
                const form = document.getElementById('check-availability-form');
                const formData = new FormData(form);
                formData.append('csrf_token', csrf_token);
                formData.append('room_id', `${room_id}`);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then((res) => res.json())
                    .then((data) => {
                        if (data.ok) {
                            attention.custom({
                                icon: "success",
                                showConfirmButton: false,
                                msg: `
                                            <p>Room is available</p>
                                            <p>
                                                <a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}" class="btn btn-primary">
                                                    Book Now!
                                                </a>
                                             </p>
                                           `.trim()
                            })
                        } else {
                            attention.error({
                                msg: "Room is not available"
                            })
                        }
                    });
            }
        },
    });
}