import type {Interceptor} from "@connectrpc/connect";
import {GenRandomStr} from "kk_kit/web/multi_lang/kk_id/kk_id";


export const authInterceptor: Interceptor = (next) => async (req) => {

    const token = getToken();
    if (token) {
        req.header.set("JwtAuthKey", token);
    }

    const traceId = GenRandomStr(10);
    if (traceId) {
        req.header.set("TraceId", traceId);
    }

    return await next(req);
};

function getToken(): string | null {
    return "ttttt2142141"

    try {
        return localStorage.getItem("auth_token");
    } catch (e) {
        return null;
    }
}