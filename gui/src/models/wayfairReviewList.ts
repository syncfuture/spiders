import { wayfairGetReviews } from '@/services/api';
import u from '@/u';
import moment from 'moment';
import { Reducer, Effect } from 'umi';

export interface IReviewListModelState {
    reviews: any[],
    totalCount: number,
    pageSize: number,
    status: string,
    sku: string,
    itemNo: string,
    fromDate: string,
}

export interface IReviewListModel {
    // namespace: 'wayfairReviewList';
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
    // namespace: 'wayfairReviewList',

    state: {
        reviews: [],
        totalCount: 0,
        pageSize: 20,
        status: "",
        sku: "",
        itemNo: "",
        fromDate: moment().add(-2, "M").format("YYYY-MM-DD"),
    },

    effects: {
        *getReviews({ _ }, { call, put, select }) {
            const state = yield select((x: any) => x["wayfairReviewList"]);
            const query = {
                sku: state.sku,
                itemNo: state.itemNo,
                fromDate: state.fromDate,
            };
            const resp = yield call(wayfairGetReviews, query);
            yield put({ type: 'setState', payload: { reviews: resp.Reviews ?? [], totalCount: resp.TotalCount } });
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
            f.setAttribute("action", u.BaseURI() + "/wayfair/reviews/export/");
            f.setAttribute("method", "post");
            f.setAttribute("target", "download");
            const skuInput = document.createElement("input");
            skuInput.setAttribute("type", "hidden");
            skuInput.setAttribute("name", "sku");
            skuInput.setAttribute("value", state.sku);
            f.append(skuInput);
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