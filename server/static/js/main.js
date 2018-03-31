
async function updateConnectionData() {
    try {
        var cons = await getConnections();
        renderConnections(cons);
    } catch (e) {
        console.log(e)
    }
}

function getConnections() {
    return fetch('/api/v1/connections')
        .then(response => {
            if (response.status !== 200) {
                throw ('Error getting connections data ' + response.status);
            }
            return response.json()
        })
}

function renderConnections(data) {
    data.forEach(c => {
        var temp = `<div class="card"><div class="card-body">
                        <div class="card-title">${c.name}</div>
                        <p>${c.server_add}</p>
                        <p>Username: ${c.username}</p>
                        <a class="btn btn-info" href="/connection/${c.id}" role="button">View</a>
                        <button type="button" class="btn btn-warning" onclick='editConnection(${JSON.stringify(c)})' >Edit</button>
                    </div></div>`
        $("#connectionCont").append(temp);
    });
}

function editConnection(c) {
    console.log(c)
    //c = JSON.parse(c)
    $("#myModal").load("connForm.html", function () {
        $('#myModal').modal({
            show: true
        });
        $("#inputID").val(c.id);
        $("#inputName").val(c.name);
        $("#inputServer").val(c.server_add);
        $("#inputUsername").val(c.username);
        $("#inputPassword").val(c.password);
    });
}

function addConnection() {
    $("#myModal").load("connForm.html", function () {
        $('#myModal').modal({
            show: true
        });
        //$("#inputID").val("");
    });
}

function submitConnForm(event) {
    pData = {}
    $('#connectionForm').serializeArray().forEach(d => {
        pData[d.name] = d.value;
    });
    if (pData.id.length < 1) {
        fetch(`/api/v1/connection`, {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'content-type': 'application/json'
                },
                body: JSON.stringify(pData),

            })
            .then(response => {
                if (response.status !== 200) {
                    throw (response);
                }
                return response.json()
            })
            .then(data => {
                console.log(data)
                $('#myModal').modal('hide');
            }) //TODO: NEED TO ADD NEW CONNECTION TO SCREEN
            .catch(e => console.log(e))
        return;
    }
    fetch(`/api/v1/connection/${pData.id}`, {
            method: 'PUT',
            credentials: 'same-origin',
            headers: {
                'content-type': 'application/json'
            },
            body: JSON.stringify(pData),
        })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            $('#myModal').modal('hide');
        })
        .catch(e => console.log(e))
}
