import { getReviews } from '@/services/api';
import moment from 'moment';
import { Reducer, Effect } from 'umi';

export interface IReviewListModelState {
    reviews: any[],
    totalCount: number,
    pageSize: number,
    status: string,
    asin: string,
    itemNo: string,
    fromDate: string,
}

export interface IReviewListModel {
    namespace: 'reviewList';
    state: IReviewListModelState;
    effects: {
        getReviews: Effect;
    };
    reducers: {
        setState: Reducer<IReviewListModelState>;
        export: Reducer<IReviewListModelState>;
    };
    // subscriptions: { setup: Subscription };
}

const ReviewListModel: IReviewListModel = {
    namespace: 'reviewList',

    state: {
        reviews: [],
        totalCount: 0,
        pageSize: 20,
        status: "",
        asin: "",
        itemNo: "",
        fromDate: moment().add(-1, "M").format("YYYY-MM-DD"),
    },

    effects: {
        *getReviews({ _ }, { call, put, select }) {
            const state = yield select((x: any) => x["reviewList"]);
            const query = {
                asin: state.asin,
                itemNo: state.itemNo,
                fromDate: state.fromDate,
            };
            const reviews = yield call(getReviews, query);
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
        export(state: any, { _ }) {
            const f = document.createElement("form");
            f.setAttribute("action", "http://localhost:7000/reviews/export/");
            f.setAttribute("method", "post");
            f.setAttribute("target", "download");
            const asinInput = document.createElement("input");
            asinInput.setAttribute("type", "hidden");
            asinInput.setAttribute("name", "asin");
            asinInput.setAttribute("value", state.asin);
            f.append(asinInput);
            const itemNoInput = document.createElement("input");
            itemNoInput.setAttribute("type", "hidden");
            itemNoInput.setAttribute("name", "itemNo");
            itemNoInput.setAttribute("value", state.itemNo);
            f.append(itemNoInput);
            const fromDateInput = document.createElement("input");
            fromDateInput.setAttribute("type", "hidden");
            fromDateInput.setAttribute("name", "fromDate");
            fromDateInput.setAttribute("value", state.fromDate);
            f.append(fromDateInput);
            document.body.append(f);
            f.submit();
            f.remove();

            return state;
        }
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