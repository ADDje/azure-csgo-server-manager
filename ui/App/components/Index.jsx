import React from 'react'
import {IndexLink} from 'react-router'
import ServerCtl from './ServerCtl/ServerCtl.jsx'
import ServerStatus from './ServerCtl/ServerStatus.jsx'

class Index extends React.Component {
    constructor(props) {
        super(props)
    }

    componentDidMount() {
        this.props.getConfigs()
        this.props.getTemplates()
    }

    render() {
        return(
            <div className="content-wrapper" style={{height: "100%"}}>
                <section className="content-header" style={{height: "100%"}}>
                <h1>
                    Azure CS:GO Server Manager
                    <small>Control your Azure CS:GO Servers</small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard" />Server Control</IndexLink></li>
                </ol>
                </section>

                <section className="content">

                <ServerStatus 
                    azureServerStatus={this.props.azureServerStatus}
                    reloadServers={this.props.reloadServers}
                    setStatus={this.props.setStatus}
                />

                <ServerCtl 
                    azureServerStatus={this.props.azureServerStatus}
                    deploymentTemplates={this.props.deploymentTemplates}
                    serverConfigs={this.props.serverConfigs}
                    getConfig={this.props.getConfig}
                />

                </section>
            </div>
        )
    }
}

export default Index
