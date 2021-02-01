import { getReviews } from '@/services/api';
import { Reducer, Effect } from 'umi';

export interface IReviewListModelState {
    reviews: any[],
}

export interface IReviewListModel {
    namespace: 'reviewList';
    state: IReviewListModelState;
    effects: {
        getReviews: Effect;
    };
    reducers: {
        setState: Reducer<IReviewListModelState>;
    };
    // subscriptions: { setup: Subscription };
}

const ReviewListModel: IReviewListModel = {
    namespace: 'reviewList',

    state: {
        reviews: [],
    },

    effects: {
        *getReviews({ _ }, { call, put }) {
            const reviews = yield call(getReviews);
            yield put({ type: 'setState', payload: { reviews } });
        },
    },
    reducers: {
        setState(state, action) {
            return {
                ...state,
                ...action.payload,
            };
        },
    },
    // subscriptions: {
    //     setup({ dispatch, history }) {
    //         return history.listen(({ pathname }) => {
    //             let selectedDB = -1;
    //             var t = pathname.match(/^\/db\/(\d+)$/);
    //             if (t !== null && t.length > 1) {
    //                 selectedDB = parseInt(t[1]);
    //             }

    //             dispatch({ type: "setState", payload: { SelectedDB: selectedDB } });
    //         });
    //     }
    // },
};

export default ReviewListModel;