import React from 'react'
import update from 'immutability-helper'
import sort from 'immutable-sort'
import ServerRow from './ServerRow.jsx'

class ServerList extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            order: "name"
        }

        this.sortServers = this.sortServers.bind(this)

        this.getOrderArrow = this.getOrderArrow.bind(this)
    }

    sortServers(a, b) {
        var reverseFactor = (this.state.order.substr(0, 1) === "-") ? -1 : 1
        var order = (reverseFactor < 0) ? this.state.order.substr(1) : this.state.order

        var hasNumSep = function(a) {
            return a.indexOf("-") > 0 || a.indexOf("_") > 0
        }

        var getNum = function(a) {
            var sep = (a.indexOf("-") > 0) ? "-" : "_"
            var parts = a.split(sep)
            var num = parts[parts.length-1]
            var iNum = parseInt(num)
            return iNum === NaN ? 0 : iNum
        }
        
        switch(order) {
            // If A has number but B doesn't, B comes first
            case "name":
                var aS = hasNumSep(a.name)
                var bS = hasNumSep(b.name)
                if (aS && bS) {
                    return (getNum(a.name) - getNum(b.name)) * reverseFactor
                }
                if (aS) {
                    return 1 * reverseFactor
                }
                if (bS) {
                    return -1 * reverseFactor
                }
                return a.name.localeCompare(b.name) * reverseFactor
            case "ip":
                var aHasIp = !(this.props.serverIps[a.name] === undefined || this.props.serverIps[a.name].loading || this.props.serverIps[a.name].error)
                var bHasIp = !(this.props.serverIps[b.name] === undefined || this.props.serverIps[b.name].loading || this.props.serverIps[b.name].error)

                if (aHasIp && bHasIp) {
                    var aa = this.props.serverIps[a.name].ip.split(".")
                    var bb = this.props.serverIps[b.name].ip.split(".")

                    var resulta = aa[0]*0x1000000 + aa[1]*0x10000 + aa[2]*0x100 + aa[3]*1
                    var resultb = bb[0]*0x1000000 + bb[1]*0x10000 + bb[2]*0x100 + bb[3]*1
                    
                    return (resulta-resultb) * reverseFactor
                }
                if (aHasIp) {
                    return 1 * reverseFactor
                }
                if (bHasIp) {
                    return -1 * reverseFactor
                }
                return 0
            case "status":

                var statusesA = a.properties.instanceView.statuses
                var statusesB = b.properties.instanceView.statuses
    
                // Not really a logical way of sorting these, just put bad things at the top by default
                var statusInfo = function(statuses) {
                    var good = true
                    for (var s in statuses) {
                        var parts = statuses[s].code.split("/")

                        switch (parts[0]) {
                            case "PowerState":
                                if (statuses[s].code.indexOf("running") === -1) {
                                    good = false
                                }
                                break
                            case "ProvisioningState":
                                if (statuses[s].code.indexOf("succeeded") === -1) {
                                    good = false
                                }
                                break
                        }
                        if (!good) {
                            break
                        }
                    }
                    return good
                }

                if (statusesA.length > 0 && statusesB.length > 0) {

                    var aGood = statusInfo(statusesA)
                    var bGood = statusInfo(statusesB)

                    if (aGood === bGood) {
                        return 0
                    }
                    
                    if (aGood) {
                        return 1 * reverseFactor
                    }
                    if (bGood) {
                        return -1 * reverseFactor
                    }
                    // Shouldn't happen?
                    return 0
                }
                if (statusesA.length > 0) {
                    return 1 * reverseFactor
                }
                if (statusesB.length > 0) {
                    return -1 * reverseFactor
                }
                return 0                        
        }
    }

    sortBy(criteria) {
        var order = (criteria === this.state.order) ? "-" + criteria : criteria
        this.setState({
            order: order
        })
    }

    getOrderArrow(criteria) {
        var reverse = (this.state.order.substr(0, 1) === "-")
        var order = (reverse) ? this.state.order.substr(1) : this.state.order

        if (order === criteria) {
            if (reverse) {
                return <i className="fa fa-fw fa-caret-down" />
            }
            return <i className="fa fa-fw fa-caret-up" />
        }
        return <i className="fa fa-fw fa-arrows-v" />
    }

    render() {

        var rows = null
        if (this.props.azureServerStatus.length > 0) {

            // TODO: Don't sort everytime we render, is inefficient
            var sortedServers = sort(this.props.azureServerStatus, this.sortServers)

            rows = sortedServers.map(function(server) {
                return(
                <ServerRow
                    key={server.name}
                    server={server}
                    ip={this.props.serverIps[server.name]}
                    reloadIp={this.props.reloadIp}
                />)
            }, this)
        } else {
            rows = <tr><td colSpan="5" className="text-center">No Servers Found</td></tr>
        }

        return (<div className="table-responsive">
                    <table className="table table-striped">
                        <thead>
                            <tr>
                                <th width="10%" />
                                <th><a onClick={this.sortBy.bind(this, "name")}>Name {this.getOrderArrow("name")}</a></th>
                                <th />
                                <th><a onClick={this.sortBy.bind(this, "ip")}>IP Address {this.getOrderArrow("ip")}</a></th>
                                <th><a onClick={this.sortBy.bind(this, "status")}>Status {this.getOrderArrow("status")}</a></th>
                            </tr>
                        </thead>
                        <tbody>
                            {rows}
                        </tbody>
                    </table>
                </div>)
    }

}

ServerList.propTypes = {
    azureServerStatus: React.PropTypes.array.isRequired,
    reloadIp: React.PropTypes.func.isRequired,
    serverIps: React.PropTypes.object.isRequired
}

export default ServerList