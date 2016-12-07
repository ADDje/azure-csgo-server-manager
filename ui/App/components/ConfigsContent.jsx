import React from 'react';
import {IndexLink} from 'react-router';
import ListConfigs from './Configs/ListConfigs.jsx';
import ConfigEditor from './Configs/ConfigEditor.jsx';

class ConfigsContent extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            listConfigs: {},
            selectedConfig: null,
            selectedConfigName: null,
        }

        this.loadConfigList = this.loadConfigList.bind(this)
        this.focusConfig = this.focusConfig.bind(this)
        this.reloadSelected = this.reloadSelected.bind(this)
    }

    componentDidMount() {
        this.loadConfigList()
    }

    reloadSelected() {
        if (this.state.selectedConfigName === null) {
            $.ajax({
                url: "/api/configs/get/" + this.state.selectedConfigName,
                dataType: "json",
                success: (data) => {
                    this.setState({selectedConfig: data.data})
                }
            })
        }
    }

    loadConfigList() {
        $.ajax({
            url: "/api/configs/list",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({listConfigs: data.data})
                } else {
                    this.setState({listConfigs: {}})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/configs/list', status, err.toString());
            }
        })
    }

    focusConfig(config, configName) {
        this.setState({
            selectedConfig: config,
            selectedConfigName: configName
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Configs
                    <small>Manage server configs</small>
                </h1>
                <ol className="breadcrumb">
                    <li>
                        <IndexLink to="/">
                            <i className="fa fa-dashboard fa-fw" />
                            Server Control
                        </IndexLink>
                    </li>
                    <li className="active">Configs</li>
                </ol>
                </section>

                <section className="content">

                    <ListConfigs
                        configs={this.state.listConfigs}
                        focusConfig={this.focusConfig}
                    />

                    <ConfigEditor
                        config={this.state.selectedConfig}
                        configName={this.state.selectedConfigName}
                        reloadSelected={this.reloadSelected}
                    />

                </section>
            </div>
        )
    }
}

export default ConfigsContent
