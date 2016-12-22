import React from 'react';

class ServerStatus extends React.Component {
    constructor(props) {
        super(props);
        this.formatServerStatus = this.formatServerStatus.bind(this)
    }

    formatServerStatus(serverStatus) {

        console.log(serverStatus)

        var status = serverStatus.properties.instanceView.statuses.filter(function(s) {
            return s.code.indexOf("PowerState") > -1
        })
        
        if (status.length === 0) {
            return <span className="label label-warning">Unknown</span>
        }

        var labelClass = "label label-" + ((status[0].code.indexOf("running") > 0) ? "success" : "danger")

        return <span className={labelClass}>{status[0].displayStatus}</span>
    }

    render() {

        
        var content = null
        var stop = null;
        if (this.props.azureServerStatus.length > 0) {
            content = this.props.azureServerStatus.map(function(server) {
                return(
                    <tr key={server.name}>
                        <td />
                        <td>{server.name}</td>
                        <td>{this.formatServerStatus(server)}</td>
                    </tr>
                )                                                  
            }, this);

            stop = (<div className="col-md-4">
                <button className="btn btn-block btn-danger" type="button" onClick={this.stopServer}><i className="fa fa-stop fa-fw" />Stop All CS:GO Servers</button>
            </div>)
        } else {
            content = <tr><td colSpan="3" className="text-center">No Servers Found</td></tr>
        }


        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Status</h3>
                </div>
                
                <div className="box-body">
                    <div className="table-responsive">
                        <table className="table table-striped">
                            <thead>
                                <tr>
                                    <th width="10%" />
                                    <th>Name</th>
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
    azureServerStatus: React.PropTypes.array.isRequired
}


export default ServerStatus
