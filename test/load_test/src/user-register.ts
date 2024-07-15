import { check } from "k6";
import { Trend } from "k6/metrics";
import { logWaitingTime, URL, logger } from "../util/config";
import { Options } from "k6/options"
import vu from "k6/execution"
import http from "k6/http";
import { faker } from "@faker-js/faker";

type ResponseRegister = {
    f_name: string
    l_name: string
}
type Response<T> = {
    code: number
    msg: string
    data: T
}

const metric = {
    createUser: new Trend("create_request_user", true),
}

const createUser = (f_name: string, l_name: string, email: string, password: string, address: string, roleID: number): Response<ResponseRegister> => {
    const url = URL
    const headers = {
        'Content-Type': 'application/json',
    }

    const payload = JSON.stringify({
        f_name: f_name,
        l_name: l_name,
        email: email,
        password: password,
        address: address,
        roleID: roleID
    })

    const res = http.post(`${url}/api/vehicles/add`, payload, { headers: headers})
    const response = res.json() as { data: ResponseRegister, msg: string }

    logWaitingTime({
        metric: metric.createUser,
        response: response,
        messageType: "Create User"
    })

    check(res, {
        'Success Post Data': (r) => r.status === 201,
        'Message Valid': (_) => response.msg === "success add user",
        'Data Valid': (_) => response.data.f_name !== ''
    })

    return {
        code: res.status,
        msg: response.msg,
        data: response.data
    }
}

export const options: Options = {
    vus: 10,
    iterations: 5,
    duration: '2s'
}

export default () => {
    const res = createUser(faker.person.firstName(), faker.person.lastName(), faker.internet.email(), faker.internet.password(), faker.person.jobArea(), faker.number.int({min: 2, max: 3}))
    logger.info(`Running iteration ${vu.vu.idInInstance} with first name: ${res.data.f_name}`)

}