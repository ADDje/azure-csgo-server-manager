import React from 'react';
import update from 'immutability-helper';

class ServerStatus extends React.Component {
    constructor(props) {
        super(props)
        this.formatServerStatus = this.formatServerStatus.bind(this)
        this.getButtons = this.getButtons.bind(this)

        // Maximum number of IP Queries to execute at once.
        this.maxAsyncIpQueries = 5

        // Server Info Delay (s)
        this.serverRefreshDelay = 10

        this.state = {
            loading: false,
            serverIps: {},
            ipQueriesInProgress: 0,
            refreshTimer: 0
        }

        this.getStatus = this.getStatus.bind(this)
        this.reloadTick = this.reloadTick.bind(this)

        this.clickStopAll = this.clickStopAll.bind(this)
        this.stopAll = this.stopAll.bind(this)

        this.clickSaveAll = this.clickSaveAll.bind(this)
        this.saveAll = this.saveAll.bind(this)
        
        this.clickStartAll = this.clickStartAll.bind(this)
        this.startAll = this.startAll.bind(this)

        this.clickTrashAll = this.clickTrashAll.bind(this)
        this.trashAll = this.trashAll.bind(this)

        this.displayIp = this.displayIp.bind(this)
        this.sendNextIpQuery = this.sendNextIpQuery.bind(this)
        this.getIpForServer = this.getIpForServer.bind(this)
    }

    componentWillMount() {
        this.getStatus();
        this.reloader = setInterval(this.reloadTick, 1000)
    }

    componentWillReceiveProps(nextProps) {
        var changes = {}
        var ipQueries = this.state.ipQueriesInProgress;
        for (var serverId in nextProps.azureServerStatus) {
            var server = nextProps.azureServerStatus[serverId]
            // If this is a new server we don't have information for...
            if (this.state.serverIps[server.name] === undefined) {
                // If we're not already executing too many requests
                if (ipQueries < this.maxAsyncIpQueries) {
                    ipQueries++
                    changes[server.name] = {$set: {
                        loading: true,
                        queued: false,
                        ip: ""
                    }}

                    this.getIpForServer(server.name);
                } else {
                    // Store the changes for later
                    changes[server.name] = {$set: {
                        loading: true,
                        queued: true,
                        ip: ""
                    }}
                }
            }
        }

        if (Object.keys(changes).length > 0) {
            this.setState({
                serverIps: update(this.state.serverIps, changes),
                ipQueriesInProgress: ipQueries
            })
        }
    }

    componentWillUnmount() {
        clearInterval(this.reloader)
    }

    reloadTick() {
        if (this.state.loading) {
            return
        }
        
        if (this.state.refreshTimer === 1) {
            this.getStatus();
        } else {
            this.setState({refreshTimer: this.state.refreshTimer - 1})
        }
    }

    getIpForServer(name) {
        $.get({
            url: "/api/servers/" + name + "/ip",
            success: function(name, data) {
                var param = {}
                param[name] = {$set: {loading: false, ip: data.data}}
                this.setState({
                    serverIps: update(this.state.serverIps, param),
                    ipQueriesInProgress: this.state.ipQueriesInProgress - 1
                })

                this.sendNextIpQuery();
            }.bind(this, name)
        })
    }

    sendNextIpQuery() {
        var numberNewQueries = this.maxAsyncIpQueries - this.state.ipQueriesInProgress
        if (numberNewQueries < 1) {
            return
        }

        var serversNeedingIp = Object.keys(this.state.serverIps).filter((s) => {
            return this.state.serverIps[s].queued
        })

        if (serversNeedingIp.length < 1) {
            return
        }

        var max = numberNewQueries
        if (x > serversNeedingIp.length) {
            max = serversNeedingIp.length
        }

        var changes = {}
        for (var x = 0; x < max; x++) {
            changes[serversNeedingIp[x]] = {$set: {
                loading: true,
                queued: false,
                ip: ""
            }}
            this.getIpForServer(serversNeedingIp[x])
        }

        this.setState({
            serverIps: update(this.state.serverIps, changes),
            ipQueriesInProgress: this.state.ipQueriesInProgress + x
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

    getStatus() {
        this.setState({loading: true})

        $.ajax({
            url: "/api/servers/getall",
            dataType: "json",
            success: (data) => {
                this.props.setStatus(data.data)
                this.setState({loading: false, refreshTimer: this.serverRefreshDelay})
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
            }
        })
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

    clickTrashAll() {
        swal({
            title: "Are you sure?",
            text: "This will delete all servers, their resources, hard drives and replays.",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete them all!",
            closeOnConfirm: true
        },
        function() {
            this.trashAll()
        }.bind(this));
    }

    trashAll() {
        $.post({
            url: "/api/server/delete",
            success: (resp) => {
                console.log("Deleted all")
            }
        })
    }

    clickStopAll() {
        swal({
            title: "Are you sure?",
            text: "This will stop all servers and deallocate their resources.",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, stop them!",
            closeOnConfirm: true
        },
        function() {
            this.stopAll()
        }.bind(this));
    }

    stopAll() {
        $.post({
            url: "/api/server/stop",
            success: (resp) => {
                console.log("Stopped all")
            }
        })
    }

    clickSaveAll() {
        swal({
            title: "Are you sure?",
            text: "This may take some time and involves lots of SSH sessions.\nAlso all VMs must have the same login details",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, save them!",
            closeOnConfirm: true
        },
        function() {
            window.setTimeout(function() {
                this.saveAll()
            }.bind(this), 1000)
        }.bind(this))
    }

    saveAll() {

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
                    text: "Please enter the VM Username",
                    type: "input",
                    showCancelButton: true,
                    confirmButtonColor: "#DD6B55"
                },
                function(username) {
                    if (username === false || username === "")
                        return false
                    
                    swal.close()

                    window.setTimeout(function() {
                        swal({
                            title: "VM Password",
                            text: "Please enter the VM Password",
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
                                url: "/api/server/save",
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
        });
        
    }
    
    clickStartAll() {
        swal({
            title: "Are you sure?",
            text: "This will start all servers and start costing ££.",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, start them!",
            closeOnConfirm: true
        },
        function(){
            this.startAll()
        }.bind(this));
    }

    startAll() {
        $.post({
            url: "/api/server/start",
            success: (resp) => {
                console.log("Started all")
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

    getGlobalButtons() {
        return (<div>
                <div className="col-md-3">
                    <button className="btn btn-block btn-danger" type="button" onClick={this.clickTrashAll}><i className="fa fa-trash fa-fw" />Delete All CS:GO Servers</button>
                </div>
                <div className="col-md-3">
                    <button className="btn btn-block btn-danger" type="button" onClick={this.clickStopAll}><i className="fa fa-stop fa-fw" />Stop All CS:GO Servers</button>
                </div>
                <div className="col-md-3">
                    <button className="btn btn-block btn-warning" type="button" onClick={this.clickSaveAll}><i className="fa fa-film fa-fw" />Save All Replays</button>
                </div>
                <div className="col-md-3">
                    <button className="btn btn-block btn-success" type="button" onClick={this.clickStartAll}><i className="fa fa-play fa-fw" />Start All CS:GO Servers</button>
                </div>
            </div>)
    }

    displayIp(server) {
        return (this.state.serverIps[server.name] === undefined || this.state.serverIps[server.name].loading) ? "Loading..." : this.state.serverIps[server.name].ip
    }

    render() {
        
        var loading = null
        if (this.state.loading) {
            loading = <i className="fa fa-refresh fa-spin" />
        } else {
            loading = <span>({this.state.refreshTimer})</span>
        }

        var content = null
        var buttons = null;
        if (this.props.azureServerStatus.length > 0) {

            content = this.props.azureServerStatus.map(function(server) {
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
            }, this);

            buttons = this.getGlobalButtons()
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
                        {buttons}
                    </div>
                </div>
            </div>
        )
    }
}

ServerStatus.propTypes = {
    azureServerStatus: React.PropTypes.array.isRequired,
    setStatus: React.PropTypes.func.isRequired
}


export default ServerStatus
