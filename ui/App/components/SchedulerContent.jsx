import React from 'react'
import {IndexLink} from 'react-router'
import ListActions from './Scheduler/ListActions.jsx'
import ScheduleEditor from './Scheduler/ScheduleEditor.jsx'
import SchedulerGuide from './Scheduler/SchedulerGuide.jsx'

class SchedulerContent extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            selectedAction: null,
            selectedActionName: ""
        }

        this.addOrUpdate = this.addOrUpdate.bind(this)
        this.selectAction = this.selectAction.bind(this)
    }

    componentDidMount() {
        this.props.getScheduleActions()
    }

    selectAction(name, action) {
        this.setState({
            selectedAction: action,
            selectedActionName: name
        })
    }

    addOrUpdate(name, action) {
        $.post({
            url: "/api/schedule/" + saveName,
            data: JSON.stringify(action, null, 4),
            dataType: "json",
            success: (data) => {
                console.log("Updated or saved: " + name)
            },
            error: (xhr, status, err) => {
                console.log('api/schedule/' + saveName, status, err.toString())
            }
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Schedule Actions
                    <small>For use with <a href="https://azure.microsoft.com/en-gb/services/scheduler/">Azure Scheduler</a></small>
                </h1>
                <ol className="breadcrumb">
                    <li><IndexLink to="/"><i className="fa fa-dashboard fa-fw" />Server Control</IndexLink></li>
                    <li className="active">Here</li>
                </ol>
                </section>

                <section className="content">
                
                    <ListActions
                        actions={this.props.scheduleActions}
                        selectedActionName={this.state.selectedActionName}
                        reloadActions={this.props.getScheduleActions}
                        selectAction={this.selectAction}
                    />

                    <ScheduleEditor
                        selectedActionName={this.state.selectedActionName}
                        actions={this.props.scheduleActions}
                        reloadSelected={this.props.getScheduleActions}
                        selectAction={this.selectAction}
                    />

                    <SchedulerGuide />

                </section>
            </div>
        )
    }
}

export default SchedulerContent
