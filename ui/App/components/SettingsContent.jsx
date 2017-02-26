import React from 'react'
import {IndexLink} from 'react-router'
import update from 'immutability-helper'

class SettingsContent extends React.Component {
    constructor(props) {
        super(props)
        this.getSettings = this.getSettings.bind(this)
        this.updateSettings = this.updateSettings.bind(this)
        this.updateClick = this.updateClick.bind(this)
        this.state = {
            settings: {}
        }
    }

    componentDidMount() {
        this.getSettings()
    }

    getSettings() {
        $.ajax({
            url: "/api/settings",
            dataType: "json",
            success: (resp) => {
                if (resp.success === true) {
                    delete resp.data.ConfFile
                    this.setState({settings: resp.data})
                }
            },
            error: (xhr, status, err) => {
                console.log('/api/settings', status, err.toString())
            }
        })
    }

    updateSettings() {
        var settingsJSON = JSON.stringify(this.state.settings)
        $.ajax({
            url: "/api/settings",
            datatype: "json",
            type: "POST",
            data: settingsJSON,
            success: (data) => {
                console.log(data)
                if (data.success === true) {
                    swal({
                        timer: 1000,
                        title: "Nice!",
                        text: "Settings Updated",
                        type: "success"
                    })
                }
            }
        })
    }

    handleSettingsChange(key, event) {
        var param = {}
        
        if (typeof(this.state.settings[key]) === "boolean")
        {
            param[key] = {$set: (event.target.value === "true") ? true : false}
        } else {
            param[key] = {$set: event.target.value}
        }

        this.setState({
            settings: update(this.state.settings, param)
        })
    }

    updateClick(e) {
        e.preventDefault()
        this.updateSettings()
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Settings
                    <small>Manage Server Settings</small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard"/>Server Control</IndexLink></li>
                    <li className="active">Here</li>
                </ol>
                </section>
                
                <section className="content">
                    <div className="box">
                        <div className="box-header">
                            <h3 className="box-title">Server Settings</h3>
                        </div>

                        <div className="box-body">
                        <div className="row">
                            <div className="col-md-10">
                                <div className="server-settings-section">
                                    <div className="table-responsive" id="settings-table">
                                        <form ref="settingsForm" className="form-horizontal" onSubmit={this.updateServerSettings}>
                                            {Object.keys(this.state.settings).map(function(key) {
                                                var setting = this.state.settings[key]
                                                var setting_key = key.replace(/_/g, " ")
                                                return(
                                                <div className="form-group" key={key}>
                                                    <label htmlFor={key} className="control-label col-md-3">{setting_key}</label>
                                                    <div className="col-md-6">
                                                        <input 
                                                            ref={key} 
                                                            id={key} 
                                                            className="form-control" 
                                                            defaultValue={setting} 
                                                            type="text" 
                                                            onChange={this.handleSettingsChange.bind(this, key)}
                                                        />
                                                    </div>
                                                </div>
                                                )
                                            }, this)}
                                            <div className="col-xs-6">
                                                <div className="form-group">
                                                    <input className="form-control btn btn-success" type="submit" value="Update Settings" onClick={this.updateClick} />
                                                </div>
                                            </div>
                                        </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                        </div>
                    </div>
                </section>
            </div>
        )
    }
}

export default SettingsContent
