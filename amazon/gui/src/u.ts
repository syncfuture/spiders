const u = {
    BaseURI: () => {
        const r = process.env.NODE_ENV === "production" ? "/api" : "http://localhost:7000/api";
        return r;
    },
}
export default u