import React from 'react';
import DynamicConfig from './DynamicConfig.jsx'

class ServerCtl extends React.Component {
    constructor(props) {
        super(props);
        this.startServer = this.startServer.bind(this);
        this.stopServer = this.stopServer.bind(this);

        this.incrementAutosave = this.incrementAutosave.bind(this);
        this.decrementAutosave = this.decrementAutosave.bind(this);

        this.incrementAutosaveSlots = this.incrementAutosaveSlots.bind(this);
        this.decrementAutosaveSlots = this.decrementAutosaveSlots.bind(this);

        this.incrementPort = this.incrementPort.bind(this);
        this.decrementPort = this.decrementPort.bind(this);
        
        this.incrementLatency = this.incrementLatency.bind(this);
        this.decrementLatency = this.decrementLatency.bind(this);

        this.toggleAllowCmd = this.toggleAllowCmd.bind(this);
        this.toggleP2P = this.toggleP2P.bind(this);
        this.toggleAutoPause = this.toggleAutoPause.bind(this);

        var selectedConfig = null;
        var selectedConfigName = null;
        if (this.props.serverConfigs != null &&
            Object.keys(this.props.serverConfigs).length > 0) {

            var selectedConfigName = Object.keys(this.props.serverConfigs)[0]
            
            selectedConfig = this.props.serverConfigs[selectedConfigName]
        }

        this.state = {
            savefile: "",
            latency: 100,
            autosaveInterval: 5,
            autosaveSlots: 10,
            port: 34197,
            disallowCmd: false,
            peer2peer: false,
            autoPause: false,

            selectedConfigName: selectedConfigName,
            selectedConfig: selectedConfig,
        }
    }

    startServer(e) {
        e.preventDefault();
        let serverSettings = {
            savefile: this.refs.savefile.value,
            latency: Number(this.refs.latency.value), 
            autosave_interval: Number(this.refs.autosaveInterval.value),
            autosave_slots: Number(this.refs.autosaveSlots.value),
            port: Number(this.refs.port.value),
            disallow_cmd: this.refs.allowCmd.checked,
            peer2peer: this.refs.p2p.checked,
            auto_pause: this.refs.autoPause.checked,
        }
        $.ajax({
            type: "POST",
            url: "/api/server/start",
            dataType: "json",
            data: JSON.stringify(serverSettings),
            success: (resp) => {
                this.props.getServStatus();
                this.props.getStatus();
                if (resp.success === true) {
                    swal("Factorio Server Started", resp.data)
                } else {
                    swal("Error", "Error starting Factorio Server", "error")
                }
            }
        })
        this.setState({
            savefile: this.refs.savefile.value,
            latency: Number(this.refs.latency.value), 
            autosaveInterval: Number(this.refs.autosaveInterval.value),
            autosaveSlots: Number(this.refs.autosaveSlots.value),
            port: Number(this.refs.port.value),
            disallowCmd: this.refs.allowCmd.checked,
            peer2peer: this.refs.p2p.checked,
            autoPause: this.refs.autoPause.checked,
        })
    }

    stopServer(e) {
        $.ajax({
            type: "GET",
            url: "/api/server/stop",
            dataType: "json",
            success: (resp) => {
                this.props.getServStatus();
                this.props.getStatus();
                console.log(resp)
                swal(resp.data)
            }
        });
        e.preventDefault();
    }

    incrementAutosave() {
        let saveInterval = this.state.autosaveInterval + 1;
        this.setState({autosaveInterval: saveInterval})
    }

    decrementAutosave() {
        let saveInterval = this.state.autosaveInterval - 1;
        this.setState({autosaveInterval: saveInterval})
    }

    incrementAutosaveSlots() {
        let saveSlots = this.state.autosaveSlots + 1;
        this.setState({autosaveSlots: saveSlots})
    }

    decrementAutosaveSlots() {
        let saveSlots = this.state.autosaveSlots - 1;
        this.setState({autosaveSlots: saveSlots})
    }

    incrementPort() {
        let port = this.state.port + 1;
        this.setState({port: port})
    }

    decrementPort() {
        let port = this.state.port - 1;
        this.setState({port: port})
    }
    
    incrementLatency() {
        let latency = this.state.latency + 1;
        this.setState({latency: latency})
    }

    decrementLatency() {
        let latency= this.state.latency- 1;
        this.setState({latency: latency})
    }

    toggleAllowCmd() {
        let cmd = !this.state.disallowCmd
        this.setState({disallowCmd: cmd})
    }

    toggleP2P() {
        let p2p = !this.state.peer2peer;
        this.setState({peer2peer: p2p})
    }

    toggleAutoPause() {
        let pause = !this.state.autoPause;
        this.setState({autoPause: pause})
    }

    componentWillReceiveProps(nextProps) {
        if (this.state.selectedConfig == null &&
            nextProps.serverConfigs != null &&
            Object.keys(nextProps.serverConfigs).length > 0) {

            var firstKey = Object.keys(nextProps.serverConfigs)[0]

            this.setState({
                selectedConfig: nextProps.serverConfigs[firstKey],
                selectedConfigName: firstKey
            })
        }
    }

    render() {
        var files = []

        for(var i in this.props.serverConfigs) {
            var config = this.props.serverConfigs[i]; 
            files.push(<option key={i} value={i}>{i}</option>);   
        }

        return(
            <div className="box">
                <div className="box-header">
                    <h3 className="box-title">Server Control</h3>
                </div>
                
                <div className="box-body">

                    <form action="" onSubmit={this.startServer}>
                        <div className="form-group">
                            <div className="row">
                                <div className="col-md-4">
                                    <button className="btn btn-block btn-success" type="submit"><i className="fa fa-play fa-fw"></i>Start Factorio Server</button>
                                </div>
                                
                                <div className="col-md-4">
                                    <button className="btn btn-block btn-danger" type="button" onClick={this.stopServer}><i className="fa fa-stop fa-fw"></i>Stop Factorio Server</button>
                                </div>
                            </div>

                            <hr />
                            <label>Select Config File</label>
                            <select ref="savefile" className="form-control" onChange={this.changeConfig}>
                                {files}
                            </select>
                        </div>

                        <div className="box box-success collapsed-box">
                            <button type="button" className="btn btn-box-tool" data-widget="collapse" disabled={this.selectedConfig}>
                                <div className="box-header with-border">
                                <i className="fa fa-plus fa-fw"></i><h4 className="box-title">Advanced</h4>
                                </div>
                            </button>
                            <div className="box-body" style={{display: "none"}}>
                                <DynamicConfig configName={this.state.selectedConfigName} config={this.state.selectedConfig} />
                            </div>
                        </div>
                    </form>
                </div>
            </div>

        )
    }
}

ServerCtl.propTypes = {
    azureServerStatus: React.PropTypes.array.isRequired,
    serverConfigs: React.PropTypes.object.isRequired,
    getConfig: React.PropTypes.func.isRequired,
    getStatus: React.PropTypes.func.isRequired,
}

export default ServerCtl
