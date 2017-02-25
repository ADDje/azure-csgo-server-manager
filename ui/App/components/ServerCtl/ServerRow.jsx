import React from 'react'

class ServerRow extends React.Component {
    constructor(props) {
        super(props)

        this.state = {

        }

        this.getButtons = this.getButtons.bind(this)
        this.displayIp = this.displayIp.bind(this)
        this.formatServerStatus = this.formatServerStatus.bind(this)
    }

    displayIp(server) {
        return (this.props.ip === undefined || this.props.ip.loading) ? "Loading..." : this.props.ip.ip
    }

    
    getButtons(serverStatus) {

        var status = serverStatus.properties.instanceView.statuses.filter(function(s) {
            return s.code.indexOf("PowerState") > -1
        })

        var startDisabled = !(status.length > 0 && status[0].code.indexOf("deallocated") > -1)
        var stopDisabled = !(status.length > 0 && status[0].code.indexOf("running") > -1)
        var trashDisabled = !(status.length > 0)
        var replayDisabled = !(status.length > 0 && status[0].code.indexOf("running") > -1)

        return (<div className="vm-buttons btn-group">
            <button className="btn btn-sm btn-primary" disabled={startDisabled} onClick={this.clickStart.bind(this, serverStatus.name)}>
                <i className="fa fa-play fa-fw" />
            </button>
            <button className="btn btn-sm btn-primary" disabled={stopDisabled} onClick={this.clickStop.bind(this, serverStatus.name)}>
                <i className="fa fa-stop fa-fw" />
            </button>
            <button className="btn btn-sm btn-primary" disabled={replayDisabled} onClick={this.clickReplay.bind(this, serverStatus.name)}>
                <i className="fa fa-film fa-fw" />
            </button>
            <button className="btn btn-sm btn-primary" disabled={trashDisabled} onClick={this.clickTrash.bind(this, serverStatus.name)}>
                <i className="fa fa-trash fa-fw" />
            </button>
        </div>)
    }

    clickReplay(name) {

        // TODO: Fix callback hell
        swal({
            title: "Replay Label",
            text: "Please enter the storage label for these replays",
            type: "input",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            placeholder: "week x"
        }, function(week) {
            if (week === false || week === "")
                return false

            swal.close();
            window.setTimeout(function() {
        
                swal({
                    title: "VM Username",
                    text: "Please enter the VM Username for '" + name + "'",
                    type: "input",
                    showCancelButton: true,
                    confirmButtonColor: "#DD6B55"
                },
                function(username){
                    if (username === false || username === "")
                        return false
                    
                    swal.close()

                    window.setTimeout(function() {
                        swal({
                            title: "VM Password",
                            text: "Please enter the VM Password for '" + name + "'",
                            type: "input",
                            inputType: "password",
                            showCancelButton: true,
                            confirmButtonColor: "#DD6B55"
                        },
                        function(password){
                            console.log("pass")
                            if (password === false || password === "")
                                return false;

                            $.post({
                                url: "/api/server/" + name + "/replay",
                                data: JSON.stringify({
                                    username: username,
                                    password: password,
                                    week: week
                                }),
                                success: (resp) => {
                                    console.log("Replay saved")
                                }
                            })

                            swal.close()
                        })
                    }, 1000);
                })
            }, 1000);
        })
    }

    clickStart(name) {
        $.post({
            url: "/api/server/" + name + "/start",
            success: (resp) => {
                console.log("Started")
            }
        })
    }

    clickStop(name) {
        $.post({
            url: "/api/server/" + name + "/stop",
            success: (resp) => {
                console.log("Stopped")
            }
        })
    }

    clickTrash(name) {
        swal({
            title: "Are you sure?",
            text: "You will not be able to recover " + name,
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: true
        },
        function(){
            this.doTrash(name)
        }.bind(this));
    }

    doTrash(name) {
        $.post({
            url: "/api/server/" + name + "/delete",
            success: (resp) => {
                console.log("Deleted")
            }
        })
    }

    formatServerStatus(serverStatus) {

        var statuses = serverStatus.properties.instanceView.statuses;
        
        if (statuses.length === 0) {
            return <span className="label label-warning">Unknown</span>
        }

        var icons = []
        for (var s in statuses) {
            var parts = statuses[s].code.split("/")

            switch (parts[0]) {
                case "PowerState":
                    var labelClass = "label label-" + ((statuses[s].code.indexOf("running") > -1) ? "success" : "danger")
                    icons.push(<span key={parts[0]} className={labelClass}>{statuses[s].displayStatus}</span>)
                    break
                case "ProvisioningState":
                    if (parts[1] !== "succeeded" && parts[2] !== "deallocated") {
                        icons.push(<span key={parts[0]} className="label label-danger">{statuses[s].displayStatus}</span>)
                    }
                    break
            }
        }

        return icons
    }

    render() {
        var server = this.props.server
        var buttons = this.getButtons(server)
        var ip = this.displayIp(server)
        return(
            <tr key={server.name}>
                <td><input type="checkbox" /></td>
                <td>{server.name}</td>
                <td>{buttons}</td>
                <td>{ip}</td>
                <td>{this.formatServerStatus(server)}</td>
            </tr>
        )
    }
}

export default ServerRow