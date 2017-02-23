import React from 'react';

class SchedulerGuide extends React.Component {

    render() {
        return (
        <div className="box box-solid">
            <div className="box-header with-border">
                <h3 className="box-title">Available Parameters</h3>
            </div>
            <div className="box-body">
                <div className="box-group" id="accordion">
                    <div className="panel box box-primary">
                        <div className="box-header with-border">
                            <h4 className="box-title">
                                <a data-toggle="collapse" data-parent="#accordion" href="#collapseOne"  className="collapsed">
                                    Deploy
                                </a>
                            </h4>
                        </div>
                        <div id="collapseOne" className="panel-collapse collapse">
                            <div className="box-body">
                                <table className="table">
                                    <tbody>
                                        <tr>
                                            <th style={{width: "80px"}} />
                                            <th style={{width: "150px"}}>Parameter</th>
                                            <th>Usage</th>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>numberOfServers</td>
                                            <td>Number of servers to deploy</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>configFile</td>
                                            <td>Game config file to use on servers</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>deploymentTemplate</td>
                                            <td>Server deployment template</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>azureServerName</td>
                                            <td>Override the template name</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>vmUsername</td>
                                            <td>Override the template vm username</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>vmPassword</td>
                                            <td>Override the template vm password</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                    <div className="panel box box-primary">
                        <div className="box-header with-border">
                            <h4 className="box-title">
                                <a data-toggle="collapse" data-parent="#accordion" href="#collapseTwo" className="collapsed">
                                    Delete
                                </a>
                            </h4>
                        </div>
                        <div id="collapseTwo" className="panel-collapse collapse">
                            <div className="box-body">
                                <table className="table">
                                    <tbody>
                                        <tr>
                                            <th style={{width: "80px"}} />
                                            <th style={{width: "150px"}}>Parameter</th>
                                            <th>Usage</th>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>serverNameTemplate</td>
                                            <td>Name match to delete servers. Use {"${n}"} for number</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>numberOfServers</td>
                                            <td>Number of servers to delete</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>startingNumber</td>
                                            <td>Number of server to start deleting from. Default: 1</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                    <div className="panel box box-primary">
                        <div className="box-header with-border">
                            <h4 className="box-title">
                                <a data-toggle="collapse" data-parent="#accordion" href="#collapseThree" className="collapsed">
                                    Start
                                </a>
                            </h4>
                        </div>
                        <div id="collapseThree" className="panel-collapse collapse">
                            <div className="box-body">
                                <table className="table">
                                    <tbody>
                                        <tr>
                                            <th style={{width: "80px"}} />
                                            <th style={{width: "150px"}}>Parameter</th>
                                            <th>Usage</th>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>serverNameTemplate</td>
                                            <td>Name match to start servers. Use {"${n}"} for number</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>numberOfServers</td>
                                            <td>Number of servers to start</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>startingNumber</td>
                                            <td>Number of server to start starting from. Default: 1</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                    <div className="panel box box-primary">
                        <div className="box-header with-border">
                            <h4 className="box-title">
                                <a data-toggle="collapse" data-parent="#accordion" href="#collapseFour" className="collapsed">
                                    Stop
                                </a>
                            </h4>
                        </div>
                        <div id="collapseFour" className="panel-collapse collapse">
                            <div className="box-body">
                                <table className="table">
                                    <tbody>
                                        <tr>
                                            <th style={{width: "80px"}} />
                                            <th style={{width: "150px"}}>Parameter</th>
                                            <th>Usage</th>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>serverNameTemplate</td>
                                            <td>Name match to stop servers. Use {"${n}"} for number</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>numberOfServers</td>
                                            <td>Number of servers to stop</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>startingNumber</td>
                                            <td>Number of server to start stopping from. Default: 1</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                    <div className="panel box box-primary">
                        <div className="box-header with-border">
                            <h4 className="box-title">
                                <a data-toggle="collapse" data-parent="#accordion" href="#collapseFive" className="collapsed">
                                    Save Replays
                                </a>
                            </h4>
                        </div>
                        <div id="collapseFive" className="panel-collapse collapse">
                            <div className="box-body">
                                <table className="table">
                                    <tbody>
                                        <tr>
                                            <th style={{width: "80px"}} />
                                            <th style={{width: "150px"}}>Parameter</th>
                                            <th>Usage</th>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>replayLabel</td>
                                            <td>Name to save the replays under. E.g "week 1"</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>serverNameTemplate</td>
                                            <td>Name match to save servers. Use {"${n}"} for number</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>numberOfServers</td>
                                            <td>Number of servers to save replays for</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-yellow">Optional</span></td>
                                            <td>startingNumber</td>
                                            <td>Number of server to start saving replays from. Default: 1</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>vmUsername</td>
                                            <td>Username to access the VMs</td>
                                        </tr>
                                        <tr>
                                            <td><span className="badge bg-green">Required</span></td>
                                            <td>vmPassword</td>
                                            <td>Password to access the VMs</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>)
    }
}

export default SchedulerGuide