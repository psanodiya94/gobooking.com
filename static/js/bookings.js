
function getAvailabilityForRoomById(id) {
    document.getElementById('check-availability-button').addEventListener('click', function () {
        let html = `
                <form id="check-availability-form" action="" method="post" class="needs-validation" novalidate>
                    <div class="form-row">
                        <div class="col">
                            <div id="reservation-dates-model" class="form-row">
                                <div class="col">
                                    <input disabled required type="text" class="form-control" id="start" name="start" placeholder="Check In">
                                </div>

                                <div class="col">
                                    <input disabled required type="text" class="form-control" id="end" name="end" placeholder="Check In">
                                </div>
                            </div>
                        </div>
                    </div>              
                </form>
            `
        attention.custom({
            text: html,
            title: "Select your dates",

            willOpen: () => {
                const elem = document.getElementById('reservation-dates-model');
                new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                });
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callback: function (result) {
                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", "{{.CSRFToken}}");
                formData.append("room_id", id);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: "success",
                                text: '<p>Room is Available!</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + '&s='
                                    + data.start_date
                                    + '&e='
                                    + data.end_date
                                    + '" class="btn btn-primary">'
                                    + 'Book Now</a></p>',
                                showCancelButton: false,
                                showConfirmButton: false,
                            })
                        } else {
                            attention.error({
                                text: "No Availability!",
                            })
                        }
                    })
            }
        });
    });
}