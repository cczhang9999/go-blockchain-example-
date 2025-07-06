#!/usr/bin/env python3
"""
JSON格式化工具
使用方法: curl ... | python3 format_json.py
"""

import json
import sys

def format_json():
    try:
        # 从标准输入读取JSON
        json_str = sys.stdin.read().strip()
        
        # 解析JSON
        data = json.loads(json_str)
        
        # 格式化输出
        formatted_json = json.dumps(data, indent=2, ensure_ascii=False)
        print(formatted_json)
        
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON: {e}")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    format_json() 