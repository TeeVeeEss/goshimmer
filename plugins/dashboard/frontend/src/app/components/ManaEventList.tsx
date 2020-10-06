import {inject, observer} from "mobx-react";
import * as React from "react";
import {Card, ListGroup} from "react-bootstrap";
import ManaStore from "app/stores/ManaStore";

interface Props {
    manaStore: ManaStore;
    title: string;
    listItems;
}

@inject("manaStore")
@observer
export default class ManaPledgeRevokeList extends React.Component<Props, any> {
    render() {
        return (
            <Card>
                <Card.Body>
                    <Card.Title>
                        {this.props.title}
                    </Card.Title>
                    <ListGroup style={{
                        fontSize: '0.75rem',
                        maxHeight: '150px',
                        overflowY: 'auto'
                    }}>
                        {this.props.listItems}
                    </ListGroup>
                </Card.Body>
            </Card>
        );
    }
}