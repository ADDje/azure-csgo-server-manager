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
            vmName: "",
            adminUserName: "",
            adminPassword: "",
            numberOfServers: 10,

            selectedConfigName: selectedConfigName,
            selectedConfig: selectedConfig,

            selectedTemplateName: selectedTemplateName,
            selectedTemplate: selectedTemplate,

            validationErrors: []
        }

        this.changeVmName = this.changeVmName.bind(this)
        this.changeAdminUserName = this.changeAdminUserName.bind(this)
        this.changeAdminPassword = this.changeAdminPassword.bind(this)
        this.changeNumberOfServers = this.changeNumberOfServers.bind(this)
        this.increaseNumberOfServers = this.increaseNumberOfServers.bind(this)
        this.decreaseNumberOfServers = this.decreaseNumberOfServers.bind(this)
        this.startServer = this.startServer.bind(this)

        this.changeConfig = this.changeConfig.bind(this)
        this.changeTemplate = this.changeTemplate.bind(this)

        this.getFormClass = this.getFormClass.bind(this)
        this.getHelp = this.getHelp.bind(this)
        this.getPlaceholder = this.getPlaceholder.bind(this)
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

    isValidPassword(p) {
        if (p.length < 6 || p.length > 72)
            return false
        
        var m1 = /[A-Z]/
        var m2 = /[a-z]/
        var m3 = /[0-9]/
        var m4 = /[\/\\#\-_!"£$%^&*()'@<>?\.:+=;|]/

        var c = 0
        if (p.match(m1))
            c++
        if (p.match(m2))
            c++
        if (p.match(m3))
            c++
        if (p.match(m4))
            c++

        return c >= 3
    }

    startServer(e) {
        e.preventDefault()
        var errors = [];

        if (this.state.vmName.length < 1
            && this.getPlaceholder("vmName") === "UNDEFINED") {
            errors.push("vmName")
        }

        if (this.state.adminUserName.length < 1
            && this.getPlaceholder("adminUserName") === "UNDEFINED") {
            errors.push("adminUserName")
        }

        if (this.state.adminPassword.length < 1) {
            var placeholder = this.getPlaceholder("adminPassword")
            if (placeholder === "UNDEFINED") {
                errors.push("adminPassword")
            } else {
                // Check the predefined value is validation
                if (!this.isValidPassword(placeholder)) {
                    errors.push("adminPassword")
                }
            }
        } else {
            if (!this.isValidPassword(this.state.adminPassword)) {
                errors.push("adminPassword")
            }
        }

        if (parseInt(this.state.numberOfServers) === NaN ||
                parseInt(this.state.numberOfServers) < 1) {
            errors.push("numberOfServers")
        }
        if (this.state.selectedConfigName.length < 1) {
            errors.push("selectedConfigName")
        }
        if (this.state.selectedTemplateName.length < 1) {
            errors.push("selectedTemplateName")
        }

        this.setState({validationErrors: errors})

        if (errors.length > 0) {
            return
        }

        let serverSettings = {
            vmName: this.state.vmName,
            adminUserName: this.state.adminUserName,
            adminPassword: this.state.adminPassword,
            numberOfServers: this.state.numberOfServers,
            configFile: this.state.selectedConfigName,
            templateFile: this.state.selectedTemplateName
        }

        $.post({
            url: "/api/server/deploy",
            data: JSON.stringify(serverSettings),
            success: (resp) => {
                
            }
        })
    }

    stopServer(e) {
        $.ajax({
            type: "GET",
            url: "/api/server/stop",
            dataType: "json",
            success: (resp) => {
                swal(resp.data)
            }
        })
        e.preventDefault()
    }

    changeVmName(e) {
        this.setState({
            vmName: e.target.value
        })
    }

    changeAdminPassword(e) {
        this.setState({
            adminPassword: e.target.value
        })
    }

    changeAdminUserName(e) {
        this.setState({
            adminUserName: e.target.value
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
        this.setState({selectedConfigName: k, selectedConfig: this.props.serverConfigs[k]})
    }

    changeTemplate(e) {
        var k = e.target.value
        this.setState({selectedTemplateName: k, selectedTemplate: this.props.deploymentTemplates[k]})
    }

    getFormClass(name) {
        var className = "form-group"
        if (this.state.validationErrors.indexOf(name) !== -1) {
            className += " has-error"
        }
        return className
    }

    getPlaceholder(name) {
        if (this.state.selectedTemplate !== null && this.state.selectedTemplate.Parameters.parameters !== undefined) {
            if (this.state.selectedTemplate.Parameters.parameters[name] === undefined) {
                return "UNDEFINED"
            } else {
                return this.state.selectedTemplate.Parameters.parameters[name].value
            }
        }
        return ""
    }

    getHelp(name) {
        if (this.state.validationErrors.indexOf(name) !== -1) {
            var msg;
            switch(name) {
                case "vmName":
                    msg = "Server name is required"
                    break
                case "adminUserName":
                    msg = "VM username is required"
                    break
                case "adminPassword":
                    msg = "The supplied password must be between 6-72 characters long"
                    msg += " and must satisfy at least 3 of password complexity requirements"
                    msg += " from the following: \r\n1) Contains an uppercase character\r\n2)"
                    msg += " Contains a lowercase character\r\n3) Contains a numeric digit\r\n4)"
                    msg += " Contains a special character"
                    break
                case "numberOfServers":
                    msg = "Must have a positive number of servers"
                    break
                case "selectedConfigName":
                    msg = "Must select a Server Config"
                    break
                case "selectedTemplateName":
                    msg = "Must select a Deployment Template"
                    break
                default:
                    console.log("Unknown validation error: " + name)
            }
            return (<span className="help-block">{msg}</span>)
        }
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

                            <h4>Quick Overrides</h4>

                            <div className={this.getFormClass("vmName")}>
                                <label>Azure Server Name</label>
                                <div className="input-group">
                                    <input type="text" className="form-control" onChange={this.changeVmName} value={this.state.vmName} placeholder={this.getPlaceholder("vmName")} />
                                </div>
                                {this.getHelp("vmName")}
                            </div>

                            <div className={this.getFormClass("adminUserName")}>
                                <label>VM Username</label>
                                <div className="input-group">
                                    <input type="text" className="form-control" onChange={this.changeAdminUserName} value={this.state.adminUserName} placeholder={this.getPlaceholder("adminUserName")} />
                                </div>
                                {this.getHelp("adminUserName")}
                            </div>

                            <div className={this.getFormClass("adminPassword")}>
                                <label>VM Password</label>
                                <div className="input-group">
                                    <input type="Password" className="form-control" onChange={this.changeAdminPassword} value={this.state.adminPassword} placeholder={this.getPlaceholder("adminPassword")} />
                                </div>
                                {this.getHelp("adminPassword")}
                            </div>

                            <h4>Server Settings</h4>
                            
                            <div className={this.getFormClass("numberOfServers")}>
                                <label>Number of Servers</label>
                                <div className="input-group">
                                    <input type="text" className="form-control" onChange={this.changeNumberOfServers} value={this.state.numberOfServers} />
                                    <div className="input-group-btn">
                                        <button type="button" className="btn btn-primary" onClick={this.increaseNumberOfServers}>
                                            <i className="fa fa-arrow-up" />
                                        </button>
                                        <button type="button" className="btn btn-primary" onClick={this.decreaseNumberOfServers}>
                                            <i className="fa fa-arrow-down" />
                                        </button>
                                    </div>
                                </div>
                                {this.getHelp("numberOfServers")}
                            </div>

                            <div className={this.getFormClass("selectedConfigName")}>
                                <label>Select Config File</label>
                                <select value={this.state.selectedConfigName} className="form-control" onChange={this.changeConfig}>
                                    {files}
                                </select>
                                {this.getHelp("selectedConfigName")}
                            </div>

                            <div className={this.getFormClass("selectedTemplateName")}>
                                <label>Select Deployment Template</label>
                                <select value={this.state.selectedTemplateName} className="form-control" onChange={this.changeTemplate}>
                                    {templates}
                                </select>
                                {this.getHelp("selectedTemplateName")}
                            </div>
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
    serverConfigs: React.PropTypes.object.isRequired,
}

export default ServerCtl
