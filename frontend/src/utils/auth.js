export const saveUserId = (id) => {
  localStorage.setItem("user_id", id);
};

export const getUserId = () => {
  return localStorage.getItem("user_id");
};

export const logoutUser = () => {
  localStorage.removeItem("user_id");
};
