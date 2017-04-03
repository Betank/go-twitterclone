import React, {Component, PropTypes} from 'react'
import {Row, Col, Grid} from 'react-bootstrap'
import {fetchStatistics} from '../actions/twitterActions'
import {connect} from 'react-redux'

class ProfileStatistics extends Component {

    componentDidMount() {
        this.props.dispatch(fetchStatistics())
    }

    render() {
        return (
            <Grid fluid>
                <Row className="show-grid">
                    <Col sm={4} md={4}>Tweets<span className="badge">{this.props.stats.tweets}</span>
                    </Col>
                    <Col sm={4} md={4}>Following<span className="badge">{this.props.stats.follow}</span>
                    </Col>
                    <Col sm={4} md={4}>Followers<span className="badge">{this.props.stats.follower}</span>
                    </Col>
                </Row>
            </Grid>
        )
    }
}

ProfileStatistics.propTypes = {
    stats: PropTypes.shape({
        follow: PropTypes.number.isRequired,
        follower: PropTypes.number.isRequired,
        tweets: PropTypes.number.isRequired
    }).isRequired
}

const mapStateToProps = (state) => {
    return {
        stats: state.stats
    }
}

export default connect(mapStateToProps)(ProfileStatistics)