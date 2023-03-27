import os
import time
import base64
import hashlib
import hmac
import json
from .miaccount import MiAccount

# REGIONS = ['cn', 'de', 'i2', 'ru', 'sg', 'us']


class MiIOService:
    def __init__(self, account: MiAccount, region=None):
        self.account = account
        self.server = (
            "https://"
            + ("" if region is None or region == "cn" else region + ".")
            + "api.io.mi.com/app"
        )

    async def miio_request(self, uri, data):
        def prepare_data(token, cookies):
            cookies["PassportDeviceId"] = token["deviceId"]
            return MiIOService.sign_data(uri, data, token["xiaomiio"][0])

        headers = {
            "User-Agent": "iOS-14.4-6.0.103-iPhone12,3--D7744744F7AF32F0544445285880DD63E47D9BE9-8816080-84A3F44E137B71AE-iPhone",
            "x-xiaomi-protocal-flag-cli": "PROTOCAL-HTTP2",
        }
        resp = await self.account.mi_request(
            "xiaomiio", self.server + uri, prepare_data, headers
        )
        if "result" not in resp:
            raise Exception(f"Error {uri}: {resp}")
        return resp["result"]

    async def home_request(self, did, method, params):
        return await self.miio_request(
            "/home/rpc/" + did,
            {
                "id": 1,
                "method": method,
                "accessKey": "IOS00026747c5acafc2",
                "params": params,
            },
        )

    async def home_get_props(self, did, props):
        return await self.home_request(did, "get_prop", props)

    async def home_set_props(self, did, props):
        return [await self.home_set_prop(did, i[0], i[1]) for i in props]

    async def home_get_prop(self, did, prop):
        return (await self.home_get_props(did, [prop]))[0]

    async def home_set_prop(self, did, prop, value):
        result = (
            await self.home_request(
                did, "set_" + prop, value if isinstance(value, list) else [value]
            )
        )[0]
        return 0 if result == "ok" else result

    async def miot_request(self, cmd, params):
        return await self.miio_request("/miotspec/" + cmd, {"params": params})

    async def miot_get_props(self, did, iids):
        params = [{"did": did, "siid": i[0], "piid": i[1]} for i in iids]
        result = await self.miot_request("prop/get", params)
        return [it.get("value") if it.get("code") == 0 else None for it in result]

    async def miot_set_props(self, did, props):
        params = [
            {"did": did, "siid": i[0], "piid": i[1], "value": i[2]} for i in props
        ]
        result = await self.miot_request("prop/set", params)
        return [it.get("code", -1) for it in result]

    async def miot_get_prop(self, did, iid):
        return (await self.miot_get_props(did, [iid]))[0]

    async def miot_set_prop(self, did, iid, value):
        return (await self.miot_set_props(did, [(iid[0], iid[1], value)]))[0]

    async def miot_action(self, did, iid, args=[]):
        result = await self.miot_request(
            "action", {"did": did, "siid": iid[0], "aiid": iid[1], "in": args}
        )
        return result.get("code", -1)

    async def device_list(self, name=None, getVirtualModel=False, getHuamiDevices=0):
        result = await self.miio_request(
            "/home/device_list",
            {
                "getVirtualModel": bool(getVirtualModel),
                "getHuamiDevices": int(getHuamiDevices),
            },
        )
        result = result["list"]
        return (
            result
            if name == "full"
            else [
                {
                    "name": i["name"],
                    "model": i["model"],
                    "did": i["did"],
                    "token": i["token"],
                }
                for i in result
                if not name or name in i["name"]
            ]
        )



    @staticmethod
    def sign_nonce(ssecurity, nonce):
        m = hashlib.sha256()
        m.update(base64.b64decode(ssecurity))
        m.update(base64.b64decode(nonce))
        return base64.b64encode(m.digest()).decode()

    @staticmethod
    def sign_data(uri, data, ssecurity):
        if not isinstance(data, str):
            data = json.dumps(data)
        nonce = base64.b64encode(
            os.urandom(8) + int(time.time() / 60).to_bytes(4, "big")
        ).decode()
        snonce = MiIOService.sign_nonce(ssecurity, nonce)
        msg = "&".join([uri, snonce, nonce, "data=" + data])
        sign = hmac.new(
            key=base64.b64decode(snonce), msg=msg.encode(), digestmod=hashlib.sha256
        ).digest()
        return {
            "_nonce": nonce,
            "data": data,
            "signature": base64.b64encode(sign).decode(),
        }
