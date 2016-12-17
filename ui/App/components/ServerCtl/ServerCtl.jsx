import React from 'react';
import swal from 'sweetalert';

class ServerCtl extends React.Component {
    constructor(props) {
        super(props)

        var selectedConfig = null
        var selectedConfigName = ""
        if (this.props.serverConfigs !== null &&
            Object.keys(this.props.serverConfigs).length > 0) {

            selectedConfigName = Object.keys(this.props.serverConfigs)[0]
            selectedConfig = this.props.serverConfigs[selectedConfigName]
        }

        var selectedTemplate = null
        var selectedTemplateName = ""
        if (this.props.deploymentTemplates !== null &&
            Object.keys(this.props.deploymentTemplates).length > 0) {
            
            selectedTemplateName = Object.keys(this.props.deploymentTemplates)[0]
            selectedTemplate = this.props.deploymentTemplates[selectedTemplateName]
        }

        this.state = {
            serverPrefix: "csgo-server-",
            serverPassword: "",
            numberOfServers: 10,

            selectedConfigName: selectedConfigName,
            selectedConfig: selectedConfig,

            selectedTemplateName: selectedTemplateName,
            selectedTemplate: selectedTemplate
        }

        this.changeServerPrefix = this.changeServerPrefix.bind(this)
        this.changeServerPassword = this.changeServerPassword.bind(this)
        this.changeNumberOfServers = this.changeNumberOfServers.bind(this)
        this.increaseNumberOfServers = this.increaseNumberOfServers.bind(this)
        this.decreaseNumberOfServers = this.decreaseNumberOfServers.bind(this)
        this.startServer = this.startServer.bind(this)

        this.changeConfig = this.changeConfig.bind(this)
        this.changeTemplate = this.changeTemplate.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if (this.state.selectedConfig === null &&
            nextProps.serverConfigs !== null &&
            Object.keys(nextProps.serverConfigs).length > 0) {

            var firstKey = Object.keys(nextProps.serverConfigs)[0]

            this.setState({
                selectedConfig: nextProps.serverConfigs[firstKey],
                selectedConfigName: firstKey
            })
        }

        if (this.state.selectedTemplate === null &&
            nextProps.deploymentTemplates !== null &&
            Object.keys(nextProps.deploymentTemplates).length > 0) {

            var firstKey = Object.keys(nextProps.deploymentTemplates)[0]

            this.setState({
                selectedTemplate: nextProps.deploymentTemplates[firstKey],
                selectedTemplateName: firstKey
            })
        }
    }

    startServer(e) {
        e.preventDefault()
        let serverSettings = {
            serverPrefix: this.state.serverPrefix,
            serverPassword: this.state.serverPassword,
            numberOfServers: this.state.numberOfServers,
            configFile: this.state.selectedConfigName,
            templateFile: this.state.selectedTemplateName
        }
        $.ajax({
            type: "POST",
            url: "/api/server/start",
            dataType: "json",
            data: serverSettings,
            success: (resp) => {
                this.props.getServStatus();
                this.props.getStatus();
                if (resp.success === true) {
                    swal("CS:GO Server Started", resp.data)
                } else {
                    swal("Error", "Error starting CS:GO Server", "error")
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
        })
        e.preventDefault()
    }

    changeServerPrefix(e) {
        this.setState({
            serverPrefix: e.target.value
        })
    }

    changeServerPassword(e) {
        this.setState({
            serverPassword: e.target.value
        })
    }

    changeNumberOfServers(e) {
        this.setState({
            numberOfServers: parseInt(e.target.value)
        })
    }

    increaseNumberOfServers() {
        this.setState({
            numberOfServers: this.state.numberOfServers + 1
        })
    }

    decreaseNumberOfServers() {
        this.setState({
            numberOfServers: this.state.numberOfServers - 1
        })
    }

    changeConfig(e) {
        var k = e.target.value
        this.setState({selectedConfigName: k, selectedConfig: [this.props.serverConfigs[k]]})
    }

    changeTemplate(e) {
        var k = e.target.value
        this.setState({selectedTemplateName: k, selectedTemplate: [this.props.deploymentTemplates[k]]})
    }

    render() {
        var files = []

        for (var i in this.props.serverConfigs) {
            var config = this.props.serverConfigs[i]
            files.push(<option key={i} value={i}>{i}</option>)
        }

        var templates = []

        for (var t in this.props.deploymentTemplates) {
            var template = this.props.deploymentTemplates[t]
            templates.push(<option key={t} value={t}>{t}</option>)
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
                                    <button className="btn btn-block btn-success" type="submit"><i className="fa fa-play fa-fw" />Start CS:GO Servers</button>
                                </div>
                            </div>

                            <hr />

                            <label>Azure Server Name Prefix</label>
                            <div className="input-group">
                                <input name="serverName" type="text" className="form-control" onChange={this.changeServerPrefix} value={this.state.serverPrefix} />
                            </div>

                            <label>Azure VM Password</label>
                            <div className="input-group">
                                <input name="vmPassword" type="Password" className="form-control" onChange={this.changeServerPassword} value={this.state.serverPassword} />
                            </div>

                            <label>Number of Servers</label>
                            <div className="input-group">
                                <input name="numberOfServers" type="text" className="form-control" onChange={this.changeNumberOfServers} value={this.state.numberOfServers} />
                                <div className="input-group-btn">
                                    <button type="button" className="btn btn-primary" onClick={this.increaseNumberOfServers}>
                                        <i className="fa fa-arrow-up" />
                                    </button>
                                    <button type="button" className="btn btn-primary" onClick={this.decreaseNumberOfServers}>
                                        <i className="fa fa-arrow-down" />
                                    </button>
                                </div>
                            </div>

                            <label>Select Config File</label>
                            <select value={this.state.selectedConfigName} className="form-control" onChange={this.changeConfig}>
                                {files}
                            </select>

                            <label>Select Deployment Template</label>
                            <select value={this.state.selectedTemplateName} className="form-control" onChange={this.changeTemplate}>
                                {templates}
                            </select>
                        </div>

                        <div className="box box-success collapsed-box">
                            <button type="button" className="btn btn-box-tool" data-widget="collapse" disabled={this.selectedConfig}>
                                <div className="box-header with-border">
                                    <i className="fa fa-plus fa-fw" /><h4 className="box-title">Advanced Server Config</h4>
                                </div>
                            </button>
                            <div className="box-body" style={{display: "none"}}>
                                {
                                    // TODO
                                }
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
    deploymentTemplates: React.PropTypes.object.isRequired,
    getConfig: React.PropTypes.func.isRequired,
    getStatus: React.PropTypes.func.isRequired,
    serverConfigs: React.PropTypes.object.isRequired,
}

export default ServerCtl
