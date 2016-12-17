import React from 'react';
import ParameterEditor from './ParameterEditor.jsx'
import TextEditor from '../Configs/TextEditor.jsx'

class TemplateViewer extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        
        if (this.props.template === null) {
            return null
        }
        

        return(
            <div className="nav-tabs-custom">
                <ul className="nav nav-tabs">
                    <li><a href="#template" data-toggle="tab">Template</a></li>
                    <li><a href="#parameter_editor" data-toggle="tab">Parameter Editor</a></li>
                </ul>
                <div className="tab-content">
                    <div className="tab-pane" id="parameter_editor">
                        <ParameterEditor
                            templateName={this.props.templateName}
                            parameters={this.props.template.Parameters}
                            reloadSelected={this.props.reloadSelected}
                        />
                    </div>
                    <div className="tab-pane" id="template">
                        <TextEditor
                            name={this.props.templateName}
                            type="template"
                            template={this.props.template.Template}
                        />
                    </div>
                </div>
            </div>
        )
    }
}

TemplateViewer.propTypes = {
    reloadSelected: React.PropTypes.func.isRequired
}

export default TemplateViewer
