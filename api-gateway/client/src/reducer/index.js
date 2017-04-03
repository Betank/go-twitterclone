import { combineReducers } from 'redux'

function tweets(state = [], action) {
    switch (action.type) {
        case "RECEIVE_TWEETS":
            return [
                ...state,
                ...action.tweets
            ]
        default:
            return state
    }
}

function user(state = {
    name: 'loading...'
}, action) {
    switch(action.type) {
        case 'RECEIVE_PROFILE':
            return Object.assign({}, state, action.user)
        default:
            return state
    }
}

function stats(state = {
    follow: 0,
    follower: 0,
    tweets: 0
}, action) {
    switch(action.type) {
        case 'RECEIVE_STATS':
            return Object.assign({}, state, action.stats)
        default:
            return state
    }
}

const twitterCloneApp = combineReducers({
    tweets,
    user,
    stats
})

export default twitterCloneApp