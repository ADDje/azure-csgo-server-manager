import React from 'react';
import {IndexLink} from 'react-router';
import ListTemplates from './Templates/ListTemplates.jsx';
import TemplateViewer from './Templates/TemplateViewer.jsx';

class TemplateContent extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            listTemplates: {},
            selectedTemplate: null,
            selectedTemplateName: null,
        }

        this.loadTemplates = this.loadTemplates.bind(this)
        this.focusTemplate = this.focusTemplate.bind(this)
    }

    componentDidMount() {
        this.loadTemplates()
    }

    loadTemplates() {
        $.ajax({
            url: "/api/templates/list",
            dataType: "json",
            success: (data) => {
                if (data.success === true) {
                    this.setState({listTemplates: data.data})
                } else {
                    this.setState({listTemplates: {}})
                }
            },
            error: (xhr, status, err) => {
                console.log('api/templates/list', status, err.toString());
            }
        })
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
                        templates={this.state.listTemplates}
                        focusTemplate={this.focusTemplate}
                    />

                    <TemplateViewer
                        template={this.state.selectedTemplate}
                        templateName={this.state.selectedTemplateName}
                    />

                </section>
            </div>
        )
    }
}

export default TemplateContent
