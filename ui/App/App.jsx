import React from 'react';
import {browserHistory} from 'react-router';
import Header from './components/Header.jsx';
import Sidebar from './components/Sidebar.jsx';
import Footer from './components/Footer.jsx';
import update from 'immutability-helper';


class App extends React.Component {
    constructor(props) {
        super(props);
        this.checkLogin = this.checkLogin.bind(this)
        this.flashMessage = this.flashMessage.bind(this)
        this.getServStatus = this.getServStatus.bind(this)
        this.getConfigs = this.getConfigs.bind(this)
        this.getConfig = this.getConfig.bind(this)
        this.getTemplates = this.getTemplates.bind(this)
        this.getStatus = this.getStatus.bind(this)
        this.getScheduleActions = this.getScheduleActions.bind(this)
        this.state = {
            serverRunning: "stopped",
            azureServerStatus: [],
            configs: {},
            templates: {},
            loggedIn: false,
            username: "",
            messages: [],
            showMessage: false,
            scheduleActions: {},
        }
    }

    componentDidMount() {
        this.checkLogin();
        // Wait 1 second before redirecting to login page
        setTimeout(() => {
            if (!this.state.loggedIn) {
                browserHistory.push("/login");
            }
        }, 1000);
    }

    flashMessage(message) {
        var m = this.state.messages;
        m.push(message);
        this.setState({messages: m, showMessage: true});
    }

    checkLogin() {
        $.ajax({
            url: "/api/user/status",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({loggedIn: true,
                        username: data.data.Username})
                }
            }
        })
    }

    getServStatus() {
        $.ajax({
            url: "/api/server/status",
            dataType: "json",
            success: (data) => {
                this.setState({serverRunning: data.data.status})
            }
        })
    }

    getConfigs() {
        $.ajax({
            url: "/api/configs/list",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({configs: data.data})
                } else {
                    this.setState({configs: []})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/configs/list', status, err.toString());
            }
        })
    }

    
    getTemplates() {
        $.ajax({
            url: "/api/templates/list",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({templates: data.data})
                } else {
                    this.setState({templates: {}})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/templates/list', status, err.toString());
            }
        })
    }

    getConfig(name) {
        $.ajax({
            url: "/api/configs/get/" + name,
            dataType: "json",
            success: (data) => {
                var config = {}

                if (data.success === true) {
                    config = data.data
                }

                var o = {}
                o[name] = {$set: config}

                var configs = update(this.state.configs, o)

                this.setState({configs: configs})
            },
            error: (xhr, status, err) => {
                console.log('api/configs/get/' + name, status, err.toString());
            }
        })
    }

    getStatus() {
        $.ajax({
            url: "/api/servers/getall",
            dataType: "json",
            success: (data) => {
                console.log(data.data)
                this.setState({azureServerStatus: data.data})
            },
            error: (xhr, status, err) => {
                console.log('api/server/status', status, err.toString());
            }
        })
    }

    getScheduleActions() {
        $.ajax({
            url: "/api/schedule/getall",
            dataType: "json",
            success: (data) => {
                console.log(data.data)
                this.setState({scheduleActions: data.data})
            }
        })
    }

    render() {
        // render main application, 
        // if logged in show application
        // if not logged in show Not logged in message
        var resp;
        if (this.state.loggedIn) {
            var resp = 
                (<div>
                    <Header 
                        username={this.state.username}
                        loggedIn={this.state.loggedIn}
                        messages={this.state.messages}
                    />

                    <Sidebar 
                        azureServerStatus={this.state.azureServerStatus}
                    />
                    
                    {// Render react-router components and pass in props
                    React.cloneElement(
                        this.props.children,
                        {message: "",
                        messages: this.state.messages,
                        flashMessage: this.flashMessage,
                        azureServerStatus: this.state.azureServerStatus,
                        getStatus: this.getStatus,
                        serverConfigs: this.state.configs,
                        deploymentTemplates: this.state.templates,
                        getConfigs: this.getConfigs,
                        getConfig: this.getConfig,
                        getTemplates: this.getTemplates,
                        username: this.state.username,
                        getServStatus: this.getServStatus,
                        reloadServers: this.getStatus,
                        getScheduleActions: this.getScheduleActions,
                        scheduleActions: this.state.scheduleActions}
                    )}

                    <Footer />
                </div>)
        } else {
            var resp = <div><p>Not Logged in</p></div>;
        }

        return(
            <div className="wrapper">
            {resp}
            </div>
        )
    }
}

export default App
