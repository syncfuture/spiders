import { Reducer, history, Subscription } from 'umi';

export interface ILayoutModelState {
    selectedPathKeys: string[],
}

export interface ILayoutModel {
    namespace: 'layout';
    state: ILayoutModelState;
    effects: {
    };
    reducers: {
        setState: Reducer<ILayoutModelState>;
        navigate: Reducer;
    };
    // subscriptions: { setup: Subscription };
}

const LayoutModel: ILayoutModel = {
    namespace: 'layout',

    state: {
        selectedPathKeys: ["/"],
    },

    effects: {
    },
    reducers: {
        setState(state, action) {
            return {
                ...state,
                ...action.payload,
            };
        },
        navigate(state: ILayoutModelState, { payload }) {
            if (state.selectedPathKeys[0] != payload.path)
            {
                state.selectedPathKeys = [payload.path];
                history.push(payload.path);
            }
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

export default LayoutModel;