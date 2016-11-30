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
        } else if (serverStatus === "stopped") {
            result = <span className="label label-danger">Not Running</span>
            return result
        } 

        return serverStatus
    }

    render() {

        var keys = Object.keys(this.props.azureServerStatus)
        
        var content = null
        var stop = null;
        if (keys.length > 0) {
            content = keys.map(function(key) {
                return(
                    <tr key={key}>
                        <td></td>
                        <td>key</td>
                        <td>{this.formatServerStatus(this.props.serverStatus[key])}</td>
                    </tr>
                )                                                  
            }, this);

            stop = <div className="col-md-4">
                <button className="btn btn-block btn-danger" type="button" onClick={this.stopServer}><i className="fa fa-stop fa-fw"></i>Stop CS:GO Servers</button>
            </div>
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
                                    <th width="10%"></th>
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
    azureServerStatus: React.PropTypes.array.isRequired,
}


export default ServerStatus
