import React from 'react';

class ServerStatus extends React.Component {
    constructor(props) {
        super(props);
        this.formatServerStatus = this.formatServerStatus.bind(this)
        this.getButtons = this.getButtons.bind(this)

        this.state = {
            loading: false
        }

        this.reload = this.reload.bind(this)
    }

    componentWillMount() {
        this.reloader = setInterval(this.reload, 5000)
    }

    componentWillUnmount() {
        clearInterval(this.reloader)
    }

    reload() {
        this.setState({loading: true});
        // TODO: Callback and do this properly
        setTimeout(function() {
            this.setState({loading: false})
        }.bind(this), 500)
        this.props.reloadServers()
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

    clickReplay(name) {
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
                        password: password
                    }),
                    success: (resp) => {
                        console.log("Replay saved")
                    }
                })

                swal.close()
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

    render() {
        
        var loading = null
        if (this.state.loading) {
            loading = <i className="fa fa-refresh fa-spin" />
        }

        var content = null
        var stop = null;
        if (this.props.azureServerStatus.length > 0) {

            content = this.props.azureServerStatus.map(function(server) {
                var buttons = this.getButtons(server)
                return(
                    <tr key={server.name}>
                        <td><input type="checkbox" /></td>
                        <td>{server.name}</td>
                        <td>{buttons}</td>
                        <td />
                        <td>{this.formatServerStatus(server)}</td>
                    </tr>
                )                                                  
            }, this);

            stop = (<div className="col-md-4">
                <button className="btn btn-block btn-danger" type="button" onClick={this.stopServer}><i className="fa fa-stop fa-fw" />Stop All CS:GO Servers</button>
            </div>)
        } else {
            content = <tr><td colSpan="5" className="text-center">No Servers Found</td></tr>
        }


        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Status</h3>
                    {loading}
                </div>
                
                <div className="box-body">
                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th width="10%" />
                                    <th>Name</th>
                                    <th />
                                    <th>IP Address</th>
                                    <th>Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                {content}
                            </tbody>
                        </table>
                        {stop}
                    </div>
                </div>
            </div>
        )
    }
}

ServerStatus.propTypes = {
    azureServerStatus: React.PropTypes.array.isRequired,
    reloadServers: React.PropTypes.func.isRequired
}


export default ServerStatus
