/*
 * @Author: Nick.nie Nick.nie@aishu.cn
 * @Date: 2022-12-14 01:30:26
 * @LastEditors: Nick.nie Nick.nie@aishu.cn
 * @LastEditTime: 2022-12-14 01:30:27
 * @FilePath: /span/config/log_config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

import "time"

var (
	Internal = 5 * time.Second
	MaxLog   = 20
)
