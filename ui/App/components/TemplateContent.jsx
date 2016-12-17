import React from 'react';
import {IndexLink} from 'react-router';
import ListTemplates from './Templates/ListTemplates.jsx';
import TemplateViewer from './Templates/TemplateViewer.jsx';

class TemplateContent extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            selectedTemplate: null,
            selectedTemplateName: null,
        }

        this.focusTemplate = this.focusTemplate.bind(this)
    }

    componentDidMount() {
        this.props.getTemplates();
    }

    focusTemplate(template, templateName) {
        this.setState({
            selectedTemplate: template,
            selectedTemplateName: templateName
        })
    }

    render() {
        return(
            <div className="content-wrapper">
                <section className="content-header">
                <h1>
                    Templates
                    <small>Manage deployment templates</small>
                </h1>
                <ol className="breadcrumb">
                    <li>
                        <IndexLink to="/">
                            <i className="fa fa-dashboard fa-fw" />
                            Server Control
                        </IndexLink>
                    </li>
                    <li className="active">Templates</li>
                </ol>
                </section>

                <section className="content">

                    <ListTemplates
                        templates={this.props.deploymentTemplates}
                        focusTemplate={this.focusTemplate}
                        reloadTemplates={this.props.getTemplates}
                    />

                    <TemplateViewer
                        template={this.state.selectedTemplate}
                        templateName={this.state.selectedTemplateName}
                        reloadSelected={this.props.getTemplates}
                    />

                </section>
            </div>
        )
    }
}

export default TemplateContent
