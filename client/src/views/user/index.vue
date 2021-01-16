<template>
  <div>
    <div class="filter-container">

      <el-select
        v-model="searchType"
        placeholder="类型"
        class="filter-item"
        style="width: 130px"
      >
        <el-option
          v-for="item in searchTypeOptions"
          :key="item.key"
          :label="item.display_name"
          :value="item.key"
        />
      </el-select>

      <el-input
        v-model="searchData"
        placeholder="内容"
        style="width: 200px;"
        class="filter-item"
        clearable
        @clear="onCleanFilter"
        @keyup.enter.native="onHandleFilter"
      />
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="onHandleFilter"
      >
        搜索
      </el-button>
      <el-button
        type="primary"
        @click="onAddUser"
      >
        增加用户
      </el-button>

    </div>

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column
        label="手机号"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.phone }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="姓名"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.real_name }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="抢购资格"
        width="150"
        align="center"
      >
        <template slot-scope="scope">

          <span v-if="scope.row.qualification == 0">无</span>
          <span v-else> 有</span>
        </template>
      </el-table-column>

      <el-table-column
        label="抢购额度"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.remain_count }}瓶</span>
        </template>
      </el-table-column>

      <el-table-column
        label="加入时间"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.create_time }}</span>
        </template>
      </el-table-column>

      <el-table-column
        width="200"
        label="操作"
        align="center"
      >
        <template slot-scope="scope">
          <span>
            <el-button
              type="primary"
              size="mini"
              @click="onSetStore(scope.row)"
            >设置商店</el-button>
          </span>
          <span>
            <el-popconfirm
              title="确定删除这个用户吗?"
              @onConfirm="onDelete(scope.row)"
            >
              <el-button
                slot="reference"
                type="danger"
                size="mini"
                icon="el-icon-delete"
                circle
              />
            </el-popconfirm>

          </span>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      layout="prev, pager, next"
      :total="totalNum"
      :page-size="20"
      @current-change="handleCurrentChange"
    />

    <el-dialog
      title="增加用户"
      :visible.sync="dialogFormVisible"
    >
      <el-form>

        <el-form-item label="手机号">
          <el-input
            v-model="phoneNumber"
            type="primary"
            placeholder="请输入手机号码"
            clearable
            style="width:300px"
          />

        </el-form-item>
        <el-form-item label="验证码">
          <el-input
            v-model="smsCode"
            type="primary"
            placeholder="请输入验证码"
            style="width:150px"
          />
          <el-button
            type="primary"
            @click="sendSms"
          >
            发送短信
          </el-button>
        </el-form-item>
      </el-form>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          type="primary"
          @click="addUser"
        >
          添加
        </el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="设置商店"
      width="80%"
      :visible.sync="storeDialogFormVisible"
    >
      <!-- <el-select
        v-model="selectStore"
        filterable
        placeholder="请选择"
      >
        <el-option
          v-for="item in storeList"
          :key="item.StoreID"
          :label=" '[' + item.StoreName + '] '+ item.StoreAddr"
          :value="item.StoreID"
        />

      </el-select> -->

      <el-table
        :data="storeList"
        @selection-change="handleStoreChange"
        style="width: 100%"
      >
        <el-table-column
          type="selection"
          width="55"
        >
        </el-table-column>
        <el-table-column
          property="StoreName"
          label="商场名"
          width="200"
        />
        <el-table-column
          property="StoreAddr"
          label="商场地址"
          width="600"
        />
      </el-table>
      <el-button
        type="primary"
        @click="onSelectStore"
      >
        确定
      </el-button>
    </el-dialog>
  </div>
</template>

<script>
import { sendSms, userAdd, userList, setStore } from '@/api/user'
import { storeList } from '@/api/store'

const searchTypeOptions = [{ key: 'phone', display_name: '手机号' }]

export default {
  name: 'Users',
  created() {
    this.fetchData()
  },
  data() {
    return {
      list: null,
      totalNum: 0,
      listLoading: false,
      listQuery: {
        page: 1,
        limit: 20,
        phone: ''
      },

      phoneNumber: '',
      smsCode: '',
      dialogFormVisible: false,
      searchTypeOptions,
      searchType: searchTypeOptions[0].key,
      searchData: '',
      storeDialogFormVisible: false,
      storeList: [],
      selectStore: null,
      currentSelectUser: null
    }
  },
  methods: {
    onAddUser() {
      this.dialogFormVisible = true
    },
    addUser() {
      userAdd({ phone: this.phoneNumber, smsCode: this.smsCode })
        .then(resp => {
          this.dialogFormVisible = false

          this.$notify({
            title: '增加成功',
            type: 'success',
            duration: 2000
          })
        })
        .catch()
    },
    sendSms() {
      sendSms({ phone: this.phoneNumber })
        .then(resp => {
          this.$notify({
            title: '短信发送成功',
            type: 'success',
            duration: 2000
          })
        })
        .catch()
    },
    fetchData() {
      this.listLoading = true
      userList(this.listQuery)
        .then(response => {
          this.list = response.users
          this.totalNum = response.totalNum
          this.listLoading = false
        })
        .catch(messge => {
          this.listLoading = false
        })
    },
    resetQuery() {
      this.listQuery = {
        page: 1,
        limit: 20,
        phone: ''
      }
    },
    onHandleFilter() {
      this.resetQuery()

      switch (this.searchType) {
        case 'phone':
          this.listQuery.phone = this.searchData
          break
      }
      this.fetchData()
    },
    onSetStore(row) {
      // 获取列表
      this.currentSelectUser = row
      storeList()
        .then(resp => {
          this.storeDialogFormVisible = true
          this.storeList = resp.list
        })
        .catch()
    },
    onCleanFilter() {
      this.resetQuery()
      this.fetchData()
    },
    handleCurrentChange(num) {
      this.listQuery.page = num
      this.fetchData()
    },
    onSelectStore() {
      if (this.currentSelectUser == null) {
        this.$message.error('没有指定用户')
        return
      }
      const storeIDs = []
      this.selectStore.forEach(element => {
        storeIDs.push(element.StoreID)
      })
      setStore({
        storeIDs: storeIDs,
        phone: this.currentSelectUser.phone
      })
        .then(resp => {
          this.$notify.success({
            title: '设置成功',
            duration: 2000
          })
          this.storeDialogFormVisible = false
          this.currentSelectUser = null
        })
        .catch()
    },
    handleStoreChange(val) {
      this.selectStore = val
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
