import React from 'react';

class ServerStatus extends React.Component {
    constructor(props) {
        super(props);
        this.formatServerStatus = this.formatServerStatus.bind(this)
    }

    formatServerStatus(serverStatus) {
        var result = {}

        if (serverStatus === "running") {
            result = <span className="label label-success">Running</span>
            return result
        } else if (serverStatus == "stopped") {
            result = <span className="label label-danger">Not Running</span>
            return result
        } 

        return serverStatus
    }

    render() {
        console.log("Server Info:");
        console.log(this.azureServerStatus);
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
                                <th>Name</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            {Object.keys(this.props.azureServerStatus).map(function(key) {
                                return(
                                    <tr key={key}>
                                        <td>key</td>
                                        <td>{this.formatServerStatus(this.props.serverStatus[key])}</td>
                                    </tr>
                                )                                                  
                            }, this)}        
                        </tbody>
                    </table>
                    </div>
                </div>
            </div>
        )
    }
}

ServerStatus.propTypes = {
    azureServerStatus: React.PropTypes.array.isRequired,
}


export default ServerStatus
