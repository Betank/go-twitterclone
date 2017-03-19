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
const twitterCloneApp = combineReducers({
    tweets,
    user
})

export default twitterCloneApp