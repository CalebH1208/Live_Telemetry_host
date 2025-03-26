import Cookies from "js-cookie";

export function setCookie(name, value, days = 7) {
  Cookies.set(name, value, { expires: days });
}

export function getCookie(name) {
  return Cookies.get(name);
}

export function removeCookie(name) {
  Cookies.remove(name);
}