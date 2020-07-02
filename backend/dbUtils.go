package main

func getIDStruct(id uint64, arrayLog []loginBan) (loginBan, int, bool) {
	// faz return do []loginban com o id pretendido, e o seu indice
	i := 0
	for i < len(arrayLog) {
		if arrayLog[i].uID == id {
			return arrayLog[i], i, true
		}
		i++
	}
	return loginBan{}, -1, false
}
func isAdmin(uid uint64) bool {
	pers := selectID(int64(uid))
	if pers.roleP == "admin" {
		return true
	}
	return false
}

func isAdminUser(user string) bool {
	pers := selectUsername(user)
	if pers.roleP == "admin" {
		return true
	}
	return false
}
