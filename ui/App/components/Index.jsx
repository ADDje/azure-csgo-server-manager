import React from 'react';
import {IndexLink} from 'react-router';
import ServerCtl from './ServerCtl/ServerCtl.jsx';
import ServerStatus from './ServerCtl/ServerStatus.jsx';

class Index extends React.Component {
    constructor(props) {
        super(props);

    }

    componentDidMount() {
        this.props.getServStatus();
        this.props.getConfigs();
        this.props.getTemplates();
        this.props.getStatus();
    }

    componentWillUnmount() {
        this.props.getServStatus();
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
                />

                <ServerCtl 
                    azureServerStatus={this.props.azureServerStatus}
                    deploymentTemplates={this.props.deploymentTemplates}
                    serverConfigs={this.props.serverConfigs}
                    getConfig={this.props.getConfig}
                    getStatus={this.props.getStatus}
                />


                </section>
            </div>
        )
    }
}

export default Index
