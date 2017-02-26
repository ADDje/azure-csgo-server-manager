import React from 'react'
import {IndexLink} from 'react-router'
import ListConfigs from './Configs/ListConfigs.jsx'
import ConfigEditor from './Configs/ConfigEditor.jsx'

class ConfigsContent extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            selectedConfig: null,
            selectedConfigName: "",
        }

        this.focusConfig = this.focusConfig.bind(this)
        this.reloadSelected = this.reloadSelected.bind(this)
    }

    componentDidMount() {
        this.props.getConfigs()
    }

    reloadSelected() {
        if (this.state.selectedConfigName !== "") {
            $.ajax({
                url: "/api/configs/get/" + this.state.selectedConfigName,
                dataType: "json",
                success: (data) => {
                    this.setState({selectedConfig: data.data})
                }
            })
        }
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
                        configs={this.props.serverConfigs}
                        focusConfig={this.focusConfig}
                        reloadConfigs={this.props.getConfigs}
                        selectedConfigName={this.state.selectedConfigName}
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
