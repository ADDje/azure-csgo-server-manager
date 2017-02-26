import React from 'react'
import update from 'immutability-helper'
import sort from 'immutable-sort'
import ServerList from './ServerList.jsx'

class ServerStatus extends React.Component {
    constructor(props) {
        super(props)

        // Maximum number of IP Queries to execute at once.
        this.maxAsyncIpQueries = 5

        // Server Info Delay (s)
        this.serverRefreshDelay = 10

        this.state = {
            serverIps: {},
            ipQueriesInProgress: 0,
            loading: false,
            refreshTimer: 0,
        }

        this.clickStopAll = this.clickStopAll.bind(this)
        this.stopAll = this.stopAll.bind(this)

        this.clickSaveAll = this.clickSaveAll.bind(this)
        this.saveAll = this.saveAll.bind(this)
        
        this.clickStartAll = this.clickStartAll.bind(this)
        this.startAll = this.startAll.bind(this)

        this.clickTrashAll = this.clickTrashAll.bind(this)
        this.trashAll = this.trashAll.bind(this)

        this.getGlobalButtons = this.getGlobalButtons.bind(this)

        this.getStatus = this.getStatus.bind(this)
        this.reloadTick = this.reloadTick.bind(this)

        this.sendNextIpQuery = this.sendNextIpQuery.bind(this)
        this.getIpForServer = this.getIpForServer.bind(this)
        this.reloadIp = this.reloadIp.bind(this)
    }

    componentWillMount() {
        this.getStatus()
        this.reloader = setInterval(this.reloadTick, 1000)
    }

    componentWillReceiveProps(nextProps) {
        var changes = {}
        var ipQueries = this.state.ipQueriesInProgress
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
                        ip: "",
                        error: false
                    }}

                    this.getIpForServer(server.name)
                } else {
                    // Store the changes for later
                    changes[server.name] = {$set: {
                        loading: true,
                        queued: true,
                        ip: "",
                        error: false
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
            this.getStatus()
        } else {
            this.setState({refreshTimer: this.state.refreshTimer - 1})
        }
    }

    reloadIp(serverName) {
        if (this.state.serverIps[serverName] === undefined) {
            return
        }

        var param = {}
        param[serverName] = {$set: {queued: true}}
        this.setState({serverIps: update(this.state.serverIps, param)}, () => {
            this.sendNextIpQuery()
        })
    }

    getIpForServer(name) {
        $.get({
            url: "/api/servers/" + name + "/ip",
            success: function(name, data) {
                var param = {}
                if (data.success) {
                    param[name] = {$set: {loading: false, ip: data.data}}
                } else {
                    param[name] = {$set: {loading: false, error: true}}
                }
                this.setState({
                    serverIps: update(this.state.serverIps, param),
                    ipQueriesInProgress: this.state.ipQueriesInProgress - 1
                }, () => {
                    this.sendNextIpQuery()
                })
            }.bind(this, name)
        })
    }

    sendNextIpQuery() {
        console.log("ip update")

        var numberNewQueries = this.maxAsyncIpQueries - this.state.ipQueriesInProgress
        if (numberNewQueries < 1) {
            return
        }

        var serversNeedingIp = Object.keys(this.state.serverIps).filter((s) => {
            return this.state.serverIps[s].queued
        })

        console.log(serversNeedingIp)

        if (serversNeedingIp.length < 1) {
            return
        }

        var max = numberNewQueries
        if (max > serversNeedingIp.length) {
            max = serversNeedingIp.length
        }

        console.log(max)

        var changes = {}
        for (var x = 0; x < max; x++) {
            changes[serversNeedingIp[x]] = {$set: {
                loading: true,
                queued: false,
                ip: "",
                error: false
            }}
            this.getIpForServer(serversNeedingIp[x])
        }

        this.setState({
            serverIps: update(this.state.serverIps, changes),
            ipQueriesInProgress: this.state.ipQueriesInProgress + x
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
        }.bind(this))
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
        }.bind(this))
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

            swal.close()
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
                                return false

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
                    }, 1000)
                })
            }, 1000)
        })
        
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
        }.bind(this))
    }

    startAll() {
        $.post({
            url: "/api/server/start",
            success: (resp) => {
                console.log("Started all")
            }
        })        
    }

    getGlobalButtons() {
        if (this.props.azureServerStatus.length < 1) {
            return null
        } 
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
    
    getStatus() {
        this.setState({loading: true})

        $.ajax({
            url: "/api/servers/getall",
            dataType: "json",
            success: (data) => {
                if (data.success) {
                    this.props.setStatus(data.data)
                    this.setState({loading: false, refreshTimer: this.serverRefreshDelay})
                } else {
                    console.error(data.data)
                }
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString())
            }
        })
    }

    render() {
        
        var loading = null
        if (this.state.loading) {
            loading = <i className="fa fa-refresh fa-spin" />
        } else {
            loading = <span>({this.state.refreshTimer})</span>
        }
        
        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Status</h3>
                    &nbsp;{loading}
                </div>
                
                <div className="box-body">
                    <ServerList
                        azureServerStatus={this.props.azureServerStatus}
                        serverIps={this.state.serverIps}
                        reloadIp={this.reloadIp}
                    />
                    {this.getGlobalButtons()}
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
