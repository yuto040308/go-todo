package password

import "golang.org/x/crypto/bcrypt"

// 生パスワード -> ハッシュ文字列（サインアップ用）
func Hash(plainText string) (string, error) {
	// 1.入力したパスワードをバイトスライスに変換
	// 2.コスト（複雑性）を引数に渡す。デフォルトは推奨値
	// 3.ハッシュ化されたパスワードをバイトスライスで受け取る
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// バイトスライスを文字列に変換して返す
	return string(hashPassword), nil
}

// ハッシュと生パスワードを照合（ログイン用）。一致した場合nil
func Verify(hash, plainText string) error {
	// パスワードを比較する
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText))
	if err != nil {
		return err
	}
	return nil
}
